package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/WAY29/version-helper/utils"

	"github.com/blang/semver"
	cli "github.com/jawher/mow.cli"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

const (
	__version__ = "4.4.0"
)

func SelfUpdatecmd(cmd *cli.Cmd) {
	var (
		assumeYes = cmd.BoolOpt("y yes", false, "Automatic yes to prompts, this means automatically updating the latest version")
	)

	cmd.Spec = "[-y | --yes]"

	cmd.Action = func() {
		var (
			input string = "y"
		)

		// 初始化日志

		latest, found, err := selfupdate.DetectLatest("WAY29/version-helper")
		if err != nil {
			utils.Error(fmt.Sprintf("%v", err), 0)
			return
		}

		v := semver.MustParse(__version__)
		if !found || latest.Version.LTE(v) {
			utils.Celebration(fmt.Sprintf("Current version-helper[%s] is the latest", __version__), 0)
			return
		}

		if !*assumeYes {
			utils.Work(fmt.Sprintf("Do you want to update pocV[%s -> %s] ? (Y/n): ", __version__, latest.Version), 0)
			input, err := bufio.NewReader(os.Stdin).ReadString('\n')
			input = strings.TrimSpace(input)
			if err != nil || (input != "y" && input != "n" && input != "") {
				utils.Error("Invalid input", 0)
				return
			}
		}

		if input == "n" {
			return
		}

		exe, err := os.Executable()
		if err != nil {
			utils.Error(err.Error(), 0)
			return
		}
		if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
			utils.Error(err.Error(), 0)
			return
		}
		utils.Celebration(fmt.Sprintf("Successfully updated to version-helper[%s]", latest.Version), 0)
	}
}
