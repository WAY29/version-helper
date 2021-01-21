# Version-Helper-3.0.0

A helper for version manage

## Usage
```
Usage: version-helper [-v] COMMAND [arg...]

A helper for version manage

Options:
  -v, --version   Show the version and exit

Commands:
  init            Create .version.toml to initialize
  set             Set version
  major           Major version upgrade
  minor           Minor version upgrade
  patch           Patch version upgrade

Run 'version-helper COMMAND --help' for more information on a command.
```

## Installation
just install it by golang 

`go get github.com/WAY29/version-helper`

or download releases from github releases

## example
```
go get github.com/WAY29/version-helper
cd project_dir/
version-helper init ; version is [major].[minor].[patch]-{banner}
; some bug fixes
version-helper patch
; some new updates
version-helper minor
; some major updates with banner
version-helper major alpha
; set banner
version-helper banner beta
; show version information
version-helepr info
```
In `.version.toml`, you can use \`\` to command execution for `[[Operate]].replece` just like bash, but does not support nesting.

## Depends
- [go-toml](https://github.com/pelletier/go-toml)
- [mow-cli](https://github.com/jawher/mow.cli)
- [goconvey](https://github.com/smartystreets/goconvey)

## Reference
[bumpversion](https://github.com/peritus/bumpversion)

## log

### V3.0.0
```
now will error when you set an invalid version. legal version is [major].[minor].[patch]. 

add serialize columns to .version.toml, And achieve backward compatibility. default pattern is "{version}-{banner}"
add test

add Depends goconvey
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

### V1.0.1
```
fix a bug that failed when operate.location is empty
fix a output when tag is true
```


