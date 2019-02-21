# Kpmenu
Kpmenu is a tool written in Go used to view a KeePass database via a dmenu, or rofi, menu.

## Features
* Supports KDBX v3.1 and v4.0 (based on [gokeepasslib](https://github.com/tobischo/gokeepasslib))
* Pretty fast database decode thanks to Go
* Interfaced with dmenu or rofi
* Customize dmenu/rofi with additional command arguments
* Kpmenu main isntance stay alive for future calls so you don't need to re-insert the password
* The open database can be cached so you don't need to re-insert the password
* Automatically put selected value into the clipboard (for a custom time) thanks to xsel
* Hidden password typing

## Dependencies
* `dmenu` or `rofi`
* `xsel`
* `go` (compile only)

## Usage
I created kpmenu to make an easy and fast way to access into my KeePass database. These are some commands that you can do:
```bash
# Open a database
kpmenu -db path/to/database.kdbx

# Open a database with a key
kpmenu -db path/to/database.kdbx -k path/to/database.key

# Open a database (credentials taken from config) with a password and rofi
kpmenu -p "mypassword" -r
```

## Installation
### From AUR
You can directly install the package [kpmenu](https://aur.archlinux.org/packages/kpmenu/).

### Compiling from source
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
Usage of kpmenu:
      --argsEntry string            Additional arguments for dmenu at entry selection, separated by a space
      --argsField string            Additional arguments for dmenu at field selection, separated by a space
      --argsMenu string             Additional arguments for dmenu at menu selection, separated by a space
      --argsPassword string         Additional arguments for dmenu at password selection, separated by a space
      --cacheOneTime                Cache the database only the first time
      --cacheTimeout int            Timeout of cache in seconds
  -c, --clip int                    Timeout of clipboard in seconds (0 = no timeout)
  -d, --database string             Path to the KeePass database
      --fieldOrder string           String order of fields to show on field selection
      --fillBlacklist string        String of blacklisted fields that won't be shown
      --fillOtherFields             Enable fill of remaining fields
  -k, --keyfile string              Path to the database keyfile
  -n, --nocache                     Disable caching of database
  -p, --password string             Password of the database
      --passwordBackground string   Color of dmenu background and text for password selection, used to hide password typing
  -r, --rofi                        Use rofi instead of dmenu
      --textEntry string            Label for entry selection
      --textField string            Label for field selection
      --textMenu string             Label for menu selection
      --textPassword string         Label for password selection
  -v, --version                     Show kpmenu version
```

## License
See the [LICENSE](https://github.com/AlessioDP/kpmenu/blob/master/LICENSE) file.
