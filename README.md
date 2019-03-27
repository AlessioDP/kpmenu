[![Go Report Card](https://goreportcard.com/badge/github.com/AlessioDP/kpmenu)](https://goreportcard.com/report/github.com/AlessioDP/kpmenu) [![Travis CI](https://travis-ci.com/AlessioDP/kpmenu.svg?branch=master)](https://travis-ci.com/AlessioDP/kpmenu)
# Kpmenu
Kpmenu is a tool written in Go used to view a KeePass database via a dmenu, or rofi, menu.

## Features
* Supports KDBX v3.1 and v4.0 (based on [gokeepasslib](https://github.com/tobischo/gokeepasslib))
* Pretty fast database decode thanks to Go
* Interfaced with dmenu or rofi
* Customize dmenu/rofi with additional command arguments
* Kpmenu main instance stay alive for future calls so you don't need to re-insert the password
* The open database can be cached so you don't need to re-insert the password
* Automatically put selected value into the clipboard (for a custom time)
  * xsel and wl-clipboard supported
  * By default it will use xsel, you can override it via config or `--clipboardTool` option
* Hidden password typing

## Dependencies
* `dmenu` or `rofi`
* `xsel` or `wl-clipboard`
* `go` (compile only)

## Usage
I created kpmenu to make an easy and fast way to access into my KeePass database. These are some commands that you can do:
```bash
# Open a database
kpmenu -d path/to/database.kdbx

# Open a database with a key
kpmenu -d path/to/database.kdbx -k path/to/database.key

# Open a database (credentials taken from config) with a password and rofi
kpmenu -p "mypassword" -r
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
```
Usage of ./kpmenu:
      --argsEntry string            Additional arguments for dmenu at entry selection, separated by a space
      --argsField string            Additional arguments for dmenu at field selection, separated by a space
      --argsMenu string             Additional arguments for dmenu at menu selection, separated by a space
      --argsPassword string         Additional arguments for dmenu at password selection, separated by a space
      --cacheOneTime                Cache the database only the first time
      --cacheTimeout int            Timeout of cache in seconds (default 60)
  -c, --clipboardTime int           Timeout of clipboard in seconds (0 = no timeout) (default 15)
      --clipboardTool string        Choose which clipboard tool to use (default "xsel")
  -d, --database string             Path to the KeePass database (default "/home/alessiodp/.kpass_pw/data.kdbx")
      --fieldOrder string           String order of fields to show on field selection (default "Password UserName URL")
      --fillBlacklist string        String of blacklisted fields that won't be shown (default "UUID")
      --fillOtherFields             Enable fill of remaining fields (default true)
  -k, --keyfile string              Path to the database keyfile (default "/home/alessiodp/.kpass_pw/data.key")
  -n, --nocache                     Disable caching of database
  -p, --password string             Password of the database
      --passwordBackground string   Color of dmenu background and text for password selection, used to hide password typing (default "black")
  -r, --rofi                        Use rofi instead of dmenu
      --textEntry string            Label for entry selection (default "Entry")
      --textField string            Label for field selection (default "Field")
      --textMenu string             Label for menu selection (default "Select")
      --textPassword string         Label for password selection (default "Password")
  -v, --version                     Show kpmenu version
```

## License
See the [LICENSE](https://github.com/AlessioDP/kpmenu/blob/master/LICENSE) file.
