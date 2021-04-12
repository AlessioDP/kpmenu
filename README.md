[![Go Report Card](https://goreportcard.com/badge/github.com/AlessioDP/kpmenu)](https://goreportcard.com/report/github.com/AlessioDP/kpmenu) [![Travis CI](https://travis-ci.com/AlessioDP/kpmenu.svg?branch=master)](https://travis-ci.com/AlessioDP/kpmenu)
# Kpmenu
Kpmenu is a tool written in Go used to view a KeePass database via a dmenu, or rofi, menu.

## Features
*   Supports KDBX v3.1 and v4.0 (based on [gokeepasslib](https://github.com/tobischo/gokeepasslib))
*   Pretty fast database decode thanks to Go
*   Interfaced with dmenu, rofi, wofi and any custom executable
*   Customize dmenu/rofi with additional command arguments
*   Kpmenu can be started as a daemon, so you don't need to re-insert credentials
    *   By default the first instance of kpmenu will enter in daemon mode (cache option) for 60 seconds
    *   You can start a permanent daemon with `--daemon` option (it won't ask open the database)
    *   Even if the cache times out, the daemon won't be killed
*   Automatically put selected value into the clipboard (for a custom time)
    *   xsel and wl-clipboard supported
    *   A custom executable can be defined for every action (copy/paste/clean clipboard)
    *   By default it will use xsel, you can override it via config or `--clipboardTool` option
    *   Hidden password typing

## Dependencies
*   `go` (compile only)

## Supports
*   `dmenu`, `rofi` and `wofi` (you can define a custom executable)
*   `xsel` and `wl-clipboard` (you can define a custom executable)

## Usage
I created kpmenu to make an easy and fast way to access into my KeePass database. These are some commands that you can do:
```bash
# Open a database
kpmenu -d path/to/database.kdbx

# Open a database with a key
kpmenu -d path/to/database.kdbx -k path/to/database.key

# Open a database (credentials taken from config) with a password and rofi
kpmenu -p "mypassword" -m rofi
```

## Installation
### From AUR
You can directly install the package [kpmenu](https://aur.archlinux.org/packages/kpmenu/).

### Compiling from source
If you do not set `$GOPATH`, go sources will be downloaded into `$HOME/go`.
```bash
# Clone repository
git clone https://github.com/AlessioDP/kpmenu
cd kpmenu

# Build
make build

# Install
sudo make install
```

## Configuration
You can set options via `config` or cli arguments.

Kpmenu will check for `$HOME/.config/kpmenu/config`, you can copy the [default one](https://github.com/AlessioDP/kpmenu/blob/master/resources/config.default) with `cp ./resources/config.default $HOME/.config/kpmenu/config`.

## Options
Options taken with `kpmenu --help`
```text
Usage of kpmenu:
      --argsEntry string              Additional arguments for dmenu at entry selection, separated by a space
      --argsField string              Additional arguments for dmenu at field selection, separated by a space
      --argsMenu string               Additional arguments for dmenu at menu selection, separated by a space
      --argsPassword string           Additional arguments for dmenu at password selection, separated by a space
      --cacheOneTime                  Cache the database only the first time
      --cacheTimeout int              Timeout of cache in seconds (default 60)
  -c, --clipboardTime int             Timeout of clipboard in seconds (0 = no timeout) (default 15)
      --clipboardTool string          Choose which clipboard tool to use (default "xsel")
      --customClipboardCopy string    Custom executable for clipboard copy
      --customClipboardPaste string   Custom executable for clipboard paste
      --customPromptEntries string    Custom executable for prompt entries
      --customPromptFields string     Custom executable for prompt fields
      --customPromptMenu string       Custom executable for prompt menu
      --customPromptPassword string   Custom executable for prompt password
      --daemon                        Start kpmenu directly as daemon
  -d, --database string               Path to the KeePass database
      --fieldOrder string             String order of fields to show on field selection (default "Password UserName URL")
      --fillBlacklist string          String of blacklisted fields that won't be shown
      --fillOtherFields               Enable fill of remaining fields (default true)
  -k, --keyfile string                Path to the database keyfile
  -m, --menu string                   Choose which menu to use (default "dmenu")
  -n, --nocache                       Disable caching of database
  -p, --password string               Password of the database
      --passwordBackground string     Color of dmenu background and text for password selection, used to hide password typing (default "black")
      --textEntry string              Label for entry selection (default "Entry")
      --textField string              Label for field selection (default "Field")
      --textMenu string               Label for menu selection (default "Select")
      --textPassword string           Label for password selection (default "Password")
  -v, --version                       Show kpmenu version
```

## License
See the [LICENSE](https://github.com/AlessioDP/kpmenu/blob/master/LICENSE) file.
