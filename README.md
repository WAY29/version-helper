# Version-Helper-4.4.1

A helper for version manage

## Usage
```
Usage: version-helper [-v] COMMAND [arg...]

A helper for version manage

Options:
  -v, --version   Show the version and exit

Commands:
  init            Create .version.toml to initialize
  info            Show version information
  set             Set version
  major           Major version upgrade
  minor           Minor version upgrade
  patch           Patch version upgrade
  banner          Set banner for version

Run 'version-helper COMMAND --help' for more information on a command.
```

## Tips
- Set environment variable LANG contains 'UTF' to enjoy beautiful unicode output!

## Installation
***make sure you can run git command before install it***

just install it by golang 

`go install github.com/WAY29/version-helper@latest`

or download releases from github releases

or build it

`go build .`

## example
```
go get github.com/WAY29/version-helper
cd project_dir/
version-helper init ; version is [major].[minor].[patch]-{banner}
# some bug fixes
version-helper patch
# some new updates
version-helper minor
# some major updates with banner
version-helper major alpha
# set banner
version-helper banner beta
# show version information
version-helepr info
# self-update
version-helper self-update
```
In `.version.toml`, you can use \`\` to command execution for `[[Operate]].replece` just like bash, but does not support nesting.

## Depends
- [go-toml](https://github.com/pelletier/go-toml)
- [mow-cli](https://github.com/jawher/mow.cli)
- [blang/semver](https://github.com/blang/semver)
- [rhysd/go-github-selfupdate](https://github.com/rhysd/go-github-selfupdate)
- [smartystreets/goconvey](https://github.com/smartystreets/goconvey)

## Reference
[bumpversion](https://github.com/peritus/bumpversion)

## log

### 4.5.0
```
add main.extraCommands for run extraCommands after update version
```

### 4.4.1
```
new tag for auto release
```
### 4.4.0
```
add goreleaser to auto release
```
### 4.3.0
```
add self-update subcommand to update version-helper
refactor code
```
### 4.2.0
```
now info subcommand will try to find .version.toml recursively from the parent directory
```
### 4.1.0
```
now version-helper will try to guarantee atomicity, it will check config before update
```
### 4.0.0
```
now tag will commit before tag, please make sure your working tree is clean. commit message is 'Update {version}'
```

### V3.0.0
```
add test for core functions
fix a bug that failed when operate.location is empty
fix a output when tag is true
```

### V2.1.0
```
Add new subcommand: info

Usage: version-helper info

Show version information
```

### V2.0.0
```
Change some output
Allow adding version number,pattern is {version}-{banner}
Add new subcommand: banner

Usage: version-helper banner [BANNER]

Set banner for version, set the banner to empty to clear the banner

Arguments:     
  BANNER       banner you want to set

```

### V1.1.0
```
Change the regular expression from non-greedy mode to greedy mode
Change some output
```

### V1.0.2
```
remove some noise success message
```