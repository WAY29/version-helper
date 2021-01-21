package cmd

import (
	"fmt"
	"strconv"

	"github.com/WAY29/version-helper/core"
	"github.com/WAY29/version-helper/utils"
	cli "github.com/jawher/mow.cli"
)

func InfoCmd(cmd *cli.Cmd) {
	cmd.Spec = ""

	cmd.Action = func() {
		config := core.LoadConfig(".version.toml")
		fmt.Println("\nVersion:", utils.WrapYellow(config.VersionHelper.Version))
		tagFlagMsg := strconv.FormatBool(config.VersionHelper.TagFlag)
		if config.VersionHelper.TagFlag {
			tagFlagMsg = utils.WrapGreen(tagFlagMsg)
		} else {
			tagFlagMsg = utils.WrapRed(tagFlagMsg)
		}
		fmt.Println("    Tag:", tagFlagMsg)
		fmt.Println("\nFiles:")
		for _, v := range config.Operate {
			fmt.Println(" - " + utils.WrapCyan(v.Location))
		}
		fmt.Println()
	}
}
