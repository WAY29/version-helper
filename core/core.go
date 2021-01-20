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
	utils.Checkf("Load "+filePath, 1)

	return Parse(data)
}

func UpgradeVersion(version string, flag int) string {
	versionSlice := strings.Split(version, ".")
	if len(versionSlice) < 3 {
		utils.Errorf("Version invalid", 1)
	}
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
	version = strings.Join(versionSlice, ".")
	utils.Checkf("Get new version", 1)
	return version
}

func Update(tomlFilePath string, config *Config) {
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
		data, err := ioutil.ReadFile(location)
		if err != nil {
			utils.Errorf("Update "+location+" : "+err.Error(), 1)
		}
		searchReg, err := regexp.Compile(strings.Replace(search, "{}", "(.*?)", -1))
		if err != nil {
			utils.Errorf("Update "+location+" : "+err.Error(), 1)
		}
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
		replace = strings.Replace(replace, "{}", version, -1)
		data = searchReg.ReplaceAll(data, []byte(replace))
		err = ioutil.WriteFile(location, data, 0666)
		if err != nil {
			utils.Errorf("Update "+location+" : "+err.Error(), 1)
		}
		utils.Checkf("Update "+location, 1)
	}
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
		utils.Checkf("git tag"+version, 1)
	}
	utils.Celebrationf("Update version to "+version+" !", 0)
}
