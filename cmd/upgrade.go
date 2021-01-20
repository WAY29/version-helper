package cmd

import (
	"github.com/WAY29/version-helper/core"
	"github.com/WAY29/version-helper/utils"
	cli "github.com/jawher/mow.cli"
)

func MajorCmd(cmd *cli.Cmd) {
	cmd.Spec = ""

	cmd.Action = func() {
		utils.Workf("Upgrade...", 0)
		config := core.Load(".version.toml")
		newVersion := core.UpgradeVersion(config.VersionHelper.Version, core.UPGRADE_MAJOR)
		config.VersionHelper.Version = newVersion
		core.Update(".version.toml", config)
	}
}

func MinorCmd(cmd *cli.Cmd) {
	cmd.Spec = ""

	cmd.Action = func() {
		utils.Workf("Upgrade...", 0)
		config := core.Load(".version.toml")
		newVersion := core.UpgradeVersion(config.VersionHelper.Version, core.UPGRADE_MINOR)
		config.VersionHelper.Version = newVersion
		core.Update(".version.toml", config)
	}
}

func PatchCmd(cmd *cli.Cmd) {
	cmd.Spec = ""

	cmd.Action = func() {
		utils.Workf("Upgrade...", 0)
		config := core.Load(".version.toml")
		newVersion := core.UpgradeVersion(config.VersionHelper.Version, core.UPGRADE_PATCH)
		config.VersionHelper.Version = newVersion
		core.Update(".version.toml", config)
	}
}

func SetCmd(cmd *cli.Cmd) {
	cmd.Spec = "VERSION"

	var (
		newVersion = cmd.StringArg("VERSION", "", "version you want to set")
	)

	cmd.Action = func() {
		utils.Workf("Upgrade...", 0)
		config := core.Load(".version.toml")
		config.VersionHelper.Version = *newVersion
		core.Update(".version.toml", config)
	}
}
