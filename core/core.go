package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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

func UpdateConfig(tomlFilePath string, oldVersion string, config *Config) {
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
		if location == "" || search == "" {
			continue
		}
		data, err := ioutil.ReadFile(location)
		if err != nil {
			utils.Errorf("Update "+location+": "+err.Error(), 1)
		}
		// ? replace search {}
		search = strings.Replace(search, "{}", oldVersion, -1)
		if err != nil {
			utils.Errorf("Update "+location+": "+err.Error(), 1)
		}
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
		// ? whether file contents has search contents
		if !bytes.Contains(data, []byte(search)) {
			utils.Errorf("Update "+location+" : "+search+" not found", 1)
		}
		// ? replace
		data = bytes.Replace(data, []byte(search), []byte(replace), -1)
		err = ioutil.WriteFile(location, data, 0666)
		if err != nil {
			utils.Errorf("Update "+location+" : "+err.Error(), 1)
		}
		utils.Checkf("Update "+location, 1)
	}
	// ? git tag
	if config.VersionHelper.TagFlag {
		stat, err := os.Stat(".git")
		if err != nil {
			utils.Warningf("Can't find .git directory", 1)
		} else if !stat.IsDir() {
			utils.Warningf(".git is not a directory", 1)
		} else {
			command := exec.Command("git", "tag", version)
			command.Start()
			err = command.Wait() //等待执行完成
			if err != nil {
				utils.Errorf("git tag "+version, 1)
			}
		}
		utils.Checkf("git tag "+version, 1)
	}
	utils.Celebrationf("Update version to "+version+" !", 0)
}
