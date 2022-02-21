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
		utils.Work("Upgrade...", 0)
		config := core.TryLoadConfig(1)
		oldVersion := config.VersionHelper.Version
		newVersion, err := core.UpgradeVersion(config, *banner, core.UPGRADE_MAJOR)
		if err != nil {
			utils.Error("Upgrade version: "+err.Error(), 1)
		}
		config.VersionHelper.Version = newVersion
		core.UpdateConfig(".version.toml", oldVersion, config)
	}
}

func MinorCmd(cmd *cli.Cmd) {
	cmd.Spec = "[BANNER]"

	var (
		banner = cmd.StringArg("BANNER", "", "Additional version number")
	)

	cmd.Action = func() {
		utils.Work("Upgrade", 0)
		config := core.TryLoadConfig(1)
		oldVersion := config.VersionHelper.Version
		newVersion, err := core.UpgradeVersion(config, *banner, core.UPGRADE_MINOR)
		if err != nil {
			utils.Error("Upgrade version: "+err.Error(), 1)
		}
		config.VersionHelper.Version = newVersion
		core.UpdateConfig(".version.toml", oldVersion, config)
	}
}

func PatchCmd(cmd *cli.Cmd) {
	cmd.Spec = "[BANNER]"

	var (
		banner = cmd.StringArg("BANNER", "", "Additional version number")
	)

	cmd.Action = func() {
		utils.Work("Upgrade version", 0)
		config := core.TryLoadConfig(1)
		oldVersion := config.VersionHelper.Version
		newVersion, err := core.UpgradeVersion(config, *banner, core.UPGRADE_PATCH)
		if err != nil {
			utils.Error("Upgrade version: "+err.Error(), 1)
		}
		config.VersionHelper.Version = newVersion
		core.UpdateConfig(".version.toml", oldVersion, config)
	}
}

func SetCmd(cmd *cli.Cmd) {
	cmd.Spec = "VERSION"

	var (
		newVersion = cmd.StringArg("VERSION", "", "version you want to set")
	)

	cmd.Action = func() {
		utils.Work("Set version", 0)
		if !core.IsPureVersion(*newVersion) {
			utils.Error("Version invalid", 1)
		}
		config := core.TryLoadConfig(1)
		serialize := config.VersionHelper.Serialize
		oldVersion := config.VersionHelper.Version
		_, banner := core.ParseVersion(serialize, oldVersion)
		temp, err := core.GenerateVersion(*newVersion, banner, config.VersionHelper.Serialize)
		if err == nil {
			*newVersion = temp
		} else {
			utils.Error("Upgrade version: "+err.Error(), 1)
		}
		config.VersionHelper.Version = *newVersion
		core.UpdateConfig(".version.toml", oldVersion, config)
	}
}

func BannerCmd(cmd *cli.Cmd) {
	cmd.Spec = "[BANNER]"
	cmd.LongDesc = "Set banner for version, set the banner to empty to clear the banner"

	var (
		banner = cmd.StringArg("BANNER", "", "banner you want to set")
	)

	cmd.Action = func() {
		utils.Work("Set banner", 0)
		config := core.TryLoadConfig(1)
		oldVersion := config.VersionHelper.Version
		newVersion, err := core.UpgradeVersion(config, *banner, core.UPGRADE_NO)
		if err != nil {
			utils.Error("Set banner: "+err.Error(), 1)
		}
		config.VersionHelper.Version = newVersion
		core.UpdateConfig(".version.toml", oldVersion, config)
	}
}
