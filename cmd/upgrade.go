package cmd

import (
	"github.com/WAY29/version-helper/core"
	"github.com/WAY29/version-helper/utils"
	cli "github.com/jawher/mow.cli"
)

func MajorCmd(cmd *cli.Cmd) {
	cmd.Spec = "[BANNER]"

	var (
		banner = cmd.StringArg("BANNER", "", "Additional version number")
	)

	cmd.Action = func() {
		utils.Workf("Upgrade...", 0)
		config := core.Load(".version.toml")
		oldVersion := config.VersionHelper.Version
		newVersion := core.UpgradeVersion(config, *banner, core.UPGRADE_MAJOR)
		config.VersionHelper.Version = newVersion
		core.Update(".version.toml", oldVersion, config)
	}
}

func MinorCmd(cmd *cli.Cmd) {
	cmd.Spec = "[BANNER]"

	var (
		banner = cmd.StringArg("BANNER", "", "Additional version number")
	)

	cmd.Action = func() {
		utils.Workf("Upgrade", 0)
		config := core.Load(".version.toml")
		oldVersion := config.VersionHelper.Version
		newVersion := core.UpgradeVersion(config, *banner, core.UPGRADE_MINOR)
		config.VersionHelper.Version = newVersion
		core.Update(".version.toml", oldVersion, config)
	}
}

func PatchCmd(cmd *cli.Cmd) {
	cmd.Spec = "[BANNER]"

	var (
		banner = cmd.StringArg("BANNER", "", "Additional version number")
	)

	cmd.Action = func() {
		utils.Workf("Upgrade", 0)
		config := core.Load(".version.toml")
		oldVersion := config.VersionHelper.Version
		newVersion := core.UpgradeVersion(config, *banner, core.UPGRADE_PATCH)
		config.VersionHelper.Version = newVersion
		core.Update(".version.toml", oldVersion, config)
	}
}

func SetCmd(cmd *cli.Cmd) {
	cmd.Spec = "VERSION"

	var (
		newVersion = cmd.StringArg("VERSION", "", "version you want to set")
	)

	cmd.Action = func() {
		utils.Workf("Upgrade", 0)
		config := core.Load(".version.toml")
		oldVersion := config.VersionHelper.Version
		_, banner := core.ParseVersion(oldVersion)
		if banner != "" {
			*newVersion += "-" + banner
		}
		config.VersionHelper.Version = *newVersion
		core.Update(".version.toml", oldVersion, config)
	}
}

func BannerCmd(cmd *cli.Cmd) {
	cmd.Spec = "[BANNER]"
	cmd.LongDesc = "Set banner for version, set the banner to empty to clear the banner"

	var (
		banner = cmd.StringArg("BANNER", "", "banner you want to set")
	)

	cmd.Action = func() {
		utils.Workf("Set banner", 0)
		config := core.Load(".version.toml")
		oldVersion := config.VersionHelper.Version
		newVersion := core.UpgradeVersion(config, *banner, core.UPGRADE_NO)
		config.VersionHelper.Version = newVersion
		core.Update(".version.toml", oldVersion, config)
	}
}
