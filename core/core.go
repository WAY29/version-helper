package core

import (
	"bytes"
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

type Config struct {
	Operate []struct {
		Location string `toml:"location"`
		Search   string `toml:"search"`
		Replace  string `toml:"replace"`
	} `toml:"operate"`
	VersionHelper struct {
		Version string `toml:"version"`
		TagFlag bool   `toml:"tag"`
	} `toml:"main"`
}

func Parse(data []byte) *Config {
	config := Config{}
	toml.Unmarshal(data, &config)
	return &config
}

func Load(filePath string) *Config {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		utils.Errorf("Load"+filePath+" : "+err.Error(), 1)
	}

	return Parse(data)
}

func ParseVersion(oldVersion string) (string, string) {
	version := oldVersion
	banner := ""
	if strings.Contains(oldVersion, "-") {
		tempStringSlice := strings.Split(oldVersion, "-")
		version = tempStringSlice[0]
		banner = tempStringSlice[1]
	}
	return version, banner
}

func UpgradeVersion(config *Config, banner string, flag int) string {
	version, _ := ParseVersion(config.VersionHelper.Version)
	// ? update cersion
	versionSlice := strings.Split(version, ".")
	if len(versionSlice) < 3 {
		utils.Errorf("Version invalid", 1)
	}
	if flag < UPGRADE_NO {
		update_int, err := strconv.Atoi(versionSlice[flag])
		if err != nil {
			utils.Errorf("Get new version : "+" : "+err.Error(), 1)
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
	// ? add banner
	version = strings.Join(versionSlice, ".")
	if banner != "" {
		version += "-" + banner
	}
	return version
}

func Update(tomlFilePath string, oldVersion string, config *Config) {
	// ? update .version.toml
	s, err := toml.Marshal(config)
	if err != nil {
		utils.Errorf("Update .version.toml"+" : "+err.Error(), 1)
	}
	s = bytes.TrimSpace(s)
	err = ioutil.WriteFile(tomlFilePath, s, 0666)
	if err != nil {
		utils.Errorf("Update .version.toml"+" : "+err.Error(), 1)
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
			utils.Errorf("Update "+location+" : "+err.Error(), 1)
		}
		// ? replace search {}
		search = strings.Replace(search, "{}", oldVersion, -1)
		if err != nil {
			utils.Errorf("Update "+location+" : "+err.Error(), 1)
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
				utils.Errorf("Update "+location+" : "+err.Error(), 1)
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
