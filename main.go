package main

import (
	"os"

	cmd "github.com/WAY29/version-helper/cmd"
	cli "github.com/jawher/mow.cli"
)

const (
	__version__ = "1.0.2"
)

func main() {
	app := cli.App("version-helper", "A helper for version manage")
	app.Version("v version", "version-helper Version: "+__version__)

	app.Spec = "[-v]"

	app.Command("init", "Create .version.toml to initialize", cmd.InitCmd)
	app.Command("set", "Set version", cmd.SetCmd)
	app.Command("major", "Major version upgrade", cmd.MajorCmd)
	app.Command("minor", "Minor version upgrade", cmd.MinorCmd)
	app.Command("patch", "Patch version upgrade", cmd.PatchCmd)

	app.Run(os.Args)
}
