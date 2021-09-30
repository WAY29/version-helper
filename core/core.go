package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/WAY29/version-helper/utils"
	toml "github.com/pelletier/go-toml"
)

const (
	UPGRADE_MAJOR = iota
	UPGRADE_MINOR
	UPGRADE_PATCH
	UPGRADE_NO
)

var EscapeStringSlice = []string{"*", ".", "?", "+", "$", "^", "[", "]", "(", ")", "|", "\\", "/"}

type Config struct {
	Operate []struct {
		Location string `toml:"location"`
		Search   string `toml:"search"`
		Replace  string `toml:"replace"`
	} `toml:"operate"`
	VersionHelper struct {
		Version   string `toml:"version"`
		TagFlag   bool   `toml:"tag"`
		Serialize string `toml:"serialize"`
	} `toml:"main"`
}

func IsPureVersion(version string) bool {
	r, _ := regexp.Compile("^\\d+\\.\\d+.\\d+$")
	if r.FindString(version) != "" {
		return true
	}
	return false
}

func parseVersionBySerialize(regString, searchString string) (resultStringMap map[string]string, err error) {
	resultStringMap = make(map[string]string)
	// ? chefk if pure version
	if IsPureVersion(searchString) {
		resultStringMap["version"] = searchString
		resultStringMap["banner"] = ""
		return
	}
	// ? get all signs
	r, _ := regexp.Compile("\\{(.*?)\\}")
	matchStringDoubleSlice := r.FindAllStringSubmatch(regString, -1)
	// ? get a new regular expression, and escape
	replacedRegString := regString
	for _, v := range EscapeStringSlice {
		replacedRegString = strings.ReplaceAll(replacedRegString, v, "\\"+v)
	}
	replacedRegString = strings.ReplaceAll(replacedRegString, "\\\\", "\\")
	replacedRegString = "^" + r.ReplaceAllString(replacedRegString, "(.*)") + "$"
	// ? check if expression is valid
	r, err = regexp.Compile(replacedRegString)
	if err != nil {
		return
	}
	// ? get result
	matchStringSlice := r.FindStringSubmatch(searchString)
	if len(matchStringSlice) < 2 {
		err = fmt.Errorf("Parse failed")
		return
	}
	for i := range matchStringSlice[1:] {
		resultStringMap[matchStringDoubleSlice[i][1]] = matchStringSlice[i+1]
	}
	return
}

func ParseVersion(serialize, oldVersion string) (string, string) {
	version := oldVersion
	banner := ""
	resultMap, err := parseVersionBySerialize(serialize, oldVersion)
	if err != nil {
		utils.Errorf("Get new verion: "+err.Error(), 1)
	}
	if resultMap["version"] != "" {
		version = resultMap["version"]
	}
	if resultMap["banner"] != "" {
		banner = resultMap["banner"]
	}
	return version, banner
}

func GenerateVersion(version, banner, serialize string) (string, error) {
	isPureVersion := IsPureVersion(version)
	// ? check version
	if !isPureVersion {
		return "", fmt.Errorf("Version invalid")
	}
	// ? check serialize
	if !strings.Contains(serialize, "{version}") || !strings.Contains(serialize, "{banner}") {
		return "", fmt.Errorf("Serialize invalid")
	}
	if isPureVersion && banner == "" {
		return version, nil
	}
	result := serialize
	result = strings.Replace(result, "{version}", version, -1)
	result = strings.Replace(result, "{banner}", banner, -1)
	return result, nil
}

func UpgradeVersion(config *Config, banner string, flag int) (string, error) {
	serialize := config.VersionHelper.Serialize
	version, _ := ParseVersion(serialize, config.VersionHelper.Version)
	// ? check version
	if !IsPureVersion(version) {
		err := fmt.Errorf("Version invalid")
		return "", err
	}
	// ? update version
	versionSlice := strings.Split(version, ".")
	if flag < UPGRADE_NO {
		update_int, err := strconv.Atoi(versionSlice[flag])
		if err != nil {
			utils.Errorf("Get new version: "+err.Error(), 1)
		}
		update_int += 1
		versionSlice[flag] = strconv.Itoa(update_int)
		if flag <= UPGRADE_MINOR {
			versionSlice[2] = "0"
		}
		if flag <= UPGRADE_MAJOR {
			versionSlice[1] = "0"
		}
	}
	// ? generate version
	version = strings.Join(versionSlice, ".")
	return GenerateVersion(version, banner, serialize)
}

func ParseConfig(data []byte) *Config {
	config := Config{}
	toml.Unmarshal(data, &config)
	// ? 2.1.0 -> 3.0.0 compatible
	if config.VersionHelper.Serialize == "" {
		config.VersionHelper.Serialize = "{version}-{banner}"
	}
	return &config
}

func LoadConfig(filePath string) *Config {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		utils.Errorf("Load"+filePath+" : "+err.Error(), 1)
	}

	return ParseConfig(data)
}

func TryLoadConfig(level int) *Config {
	configPath := FindConfigDir()
	if configPath == "" {
		utils.Errorf("Not in version-helper project directory", level)
	}
	os.Chdir(filepath.Dir(configPath))
	config := LoadConfig(configPath)
	return config
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func FindConfigDir() (configPath string) {
	configPath = ""
	currentDir, err := filepath.Abs(".")
	if err != nil {
		return
	}
	for currentDir != filepath.Dir(currentDir) {
		configPath = filepath.Join(currentDir, ".version.toml")
		if PathExists(configPath) {
			return
		}
		currentDir = filepath.Dir(currentDir)
	}
	return ""
}

func UpdateConfig(tomlFilePath string, oldVersion string, config *Config) {
	// ? check working tree if tagFlag is true
	var wtResult bool
	var err error
	var contentsMap = make(map[string][]byte, len(config.Operate))

	if config.VersionHelper.TagFlag {
		wtResult, err = CheckGit()
		if err != nil {
			utils.Errorf(err.Error(), 1)
		}
		utils.Checkf("Working tree clean", 1)
	}
	// ? guarantee atomicity
	for i, v := range config.Operate {
		location, search, _ := v.Location, v.Search, v.Replace
		search = strings.Replace(search, "{}", oldVersion, -1)
		if location == "" {
			utils.Errorf("Check config: ["+strconv.Itoa(i)+"] location invalid", 1)
		} else if search == "" {
			utils.Errorf("Check config: ["+strconv.Itoa(i)+"] search invalid", 1)
		}
		// ? read file
		content, err := ioutil.ReadFile(location)
		if err != nil {
			utils.Errorf("Check config: "+err.Error(), 1)
		}
		// ? whether file contents has search contents
		if !bytes.Contains(content, []byte(search)) {
			utils.Errorf("Check config: ["+strconv.Itoa(i)+"] "+search+" not found", 1)
		}
		contentsMap[location] = content
	}

	// ? update .version.toml
	s, err := toml.Marshal(config)
	if err != nil {
		utils.Errorf("Update .version.toml: "+err.Error(), 1)
	}
	s = bytes.TrimSpace(s)
	err = ioutil.WriteFile(tomlFilePath, s, 0666)
	if err != nil {
		utils.Errorf("Update .version.toml: "+err.Error(), 1)
	}
	utils.Checkf("Update .version.toml", 1)

	version := config.VersionHelper.Version

	// ? update action
	for _, v := range config.Operate {
		location, search, replace := v.Location, v.Search, v.Replace
		search = strings.Replace(search, "{}", oldVersion, -1)
		content := contentsMap[location]
		// ? relace `` to execute command
		shellCommandFlagCount := strings.Count(replace, "`")
		if shellCommandFlagCount > 0 && shellCommandFlagCount%2 == 0 {
			reg, err := regexp.Compile("`(.*?)`")
			if err != nil {
				utils.Errorf("Update "+location+" : "+err.Error(), 1)
			}
			commandString := reg.FindStringSubmatch(replace)[1]
			args := strings.Split(commandString, " ")
			command := exec.Command(args[0], args[1:]...)
			output, err := command.Output()
			if err != nil {
				utils.Errorf("Update "+location+": "+err.Error(), 1)
			}
			outputString := strings.TrimSpace(string(output[:]))
			replace = reg.ReplaceAllString(replace, outputString)
		}
		// ? replace replace {}
		replace = strings.Replace(replace, "{}", version, -1)
		// ? replace
		content = bytes.Replace(content, []byte(search), []byte(replace), -1)
		err = ioutil.WriteFile(location, content, 0666)
		if err != nil {
			utils.Errorf("Update "+location+" : "+err.Error(), 1)
		}
		utils.Checkf("Update "+location, 1)
	}
	// ? git commit and tag
	if wtResult && config.VersionHelper.TagFlag {
		GitCommit(config)
		GitTag(version)
	}
	utils.Celebrationf("Update version to "+version+" !", 0)
}

func CheckGit() (result bool, err error) {
	result = false
	stat, err := os.Stat(".git")
	if err != nil {
		return
	} else if !stat.IsDir() {
		err = fmt.Errorf(".git is not a directory")
	}
	command := exec.Command("git", "status")
	output, err := command.Output()
	if err == nil {
		if strings.Contains(string(output[:]), "nothing to commit, working tree clean") {
			result = true
		} else {
			err = fmt.Errorf("working tree isn't clean / git command not found")
		}
	}
	return
}

func GitCommit(config *Config) {
	defer func() {
		if err := recover(); err != nil {
			utils.Errorf(err.(error).Error(), 2)
		}
	}()
	utils.Workf("git commit", 1)
	commandArgs := make([]string, len(config.Operate)+2)
	commandArgs[0] = "add"
	commandArgs[1] = ".version.toml"
	version := config.VersionHelper.Version
	for i, v := range config.Operate {
		commandArgs[i+2] = v.Location
	}
	commandString := fmt.Sprintf("%s %s", "git", strings.Join(commandArgs, " "))
	command := exec.Command("git", commandArgs...)
	err := command.Run()
	if err != nil {
		panic(fmt.Errorf("%w (Call Command %s)", err, commandString))
	}
	msg := fmt.Sprintf("Update " + version)
	commandString = fmt.Sprintf("%s %s %s", "git", "commit", "-m "+msg)
	command = exec.Command("git", "commit", "-m "+msg)
	err = command.Run()
	if err != nil {
		panic(fmt.Errorf("%w (Call Command %s)", err, commandString))
	}
	utils.Checkf("add && commit", 2)
}

func GitTag(version string) {
	utils.Workf("git tag", 1)
	command := exec.Command("git", "tag", "v"+version)
	command.Start()
	err := command.Wait()
	if err != nil {
		utils.Errorf(err.Error(), 2)
	}
	utils.Checkf("git tag "+version, 2)
}
