# Version-Helper

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

`go get https://github.com/WAY29/version-helper`

or download releases from github releases

## example
```
go get https://github.com/WAY29/version-helper
cd project_dir/
version-helper init ; version is [major].[minor].[patch]
; some bug fixes
version-helper patch
; some new updates
version-helper minor
; some major updates
```

## Depends
- [go-toml](https://github.com/pelletier/go-toml)
- [mow-cli](https://github.com/jawher/mow.cli)