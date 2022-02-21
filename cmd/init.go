package cmd

import (
	"io/ioutil"
	"strings"

	"github.com/WAY29/version-helper/utils"
	cli "github.com/jawher/mow.cli"
)

const tomlTemplate = `[main]
version       = "{{initVersion}}"
tag           = true
extraCommands = []
serialize     = "{version}-{banner}"
[[operate]]
location = "pyproject.toml"
search   = "version = \"{}\""
replace  = "version = \"{}\""
`

func InitCmd(cmd *cli.Cmd) {
	cmd.Spec = "[VERSION]"

	var (
		initVersion = cmd.StringArg("VERSION", "0.0.1", "Initial version")
	)
	cmd.Action = func() {
		s := strings.Replace(tomlTemplate, "{{initVersion}}", *initVersion, 1)
		utils.Work("Creating .version.toml...", 0)
		err := ioutil.WriteFile(".version.toml", []byte(s), 0666)
		if err != nil {
			utils.Error("Create .version.toml", 1)
		}
		utils.Check("Create .version.toml", 1)
		utils.Celebration("Init finish, enjot it!", 0)
	}
}
