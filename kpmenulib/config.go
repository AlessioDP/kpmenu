package kpmenulib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Configuration is the main structure of kpmenu config
type Configuration struct {
	General    ConfigurationGeneral
	Executable ConfigurationExecutable
	Style      ConfigurationStyle
	Database   ConfigurationDatabase
	Flags      Flags
}

// ConfigurationGeneral is the sub-structure of the configuration related to general kpmenu settings
type ConfigurationGeneral struct {
	Menu             string // Which menu to use
	ClipboardTool    string // Clipboard tool to use
	ClipboardTimeout int    // Clipboard timeout before clean it
	NoCache          bool   // Flag to do not cache master password
	CacheOneTime     bool   // Cache the password only the first time you write it
	CacheTimeout     int    // Timeout of cache
}

// ConfigurationExecutable is the sub-structure of the configuration related to tools executed by kpmenu
type ConfigurationExecutable struct {
	CustomPromptPassword string // Custom executable for prompt password
	CustomPromptMenu     string // Custom executable for prompt menu
	CustomPromptEntries  string // Custom executable for prompt entries
	CustomPromptFields   string // Custom executable for prompt fields
	CustomClipboardCopy  string // Custom executable for clipboard copy
	CustomClipboardPaste string // Custom executable for clipboard paste
	CustomClipboardClean string // Custom executable for clipboard clean
}

// ConfigurationStyle is the sub-structure of the configuration related to style of dmenu
type ConfigurationStyle struct {
	PasswordBackground string
	TextPassword       string
	TextMenu           string
	TextEntry          string
	TextField          string
	FormatEntry        string
	ArgsPassword       string
	ArgsMenu           string
	ArgsEntry          string
	ArgsField          string
}

// ConfigurationDatabase is the sub-structure of the configuration related to database settings
type ConfigurationDatabase struct {
	Database        string
	KeyFile         string
	Password        string
	FieldOrder      string
	FillOtherFields bool
	FillBlacklist   string
}

// Flags is the sub-structure of the configuration used to handle flags that aren't into the config file
type Flags struct {
	Daemon  bool
	Version bool
}

// Menu tools used for prompts
const (
	PromptDmenu  = "dmenu"
	PromptRofi   = "rofi"
	PromptWofi   = "wofi"
	PromptCustom = "custom"
)

// Clipboard tools used for clipboard manager
const (
	ClipboardToolXsel        = "xsel"
	ClipboardToolWlclipboard = "wl-clipboard"
	ClipboardToolCustom      = "custom"
)

// NewConfiguration initializes a new Configuration pointer
func NewConfiguration() *Configuration {
	return &Configuration{
		General: ConfigurationGeneral{
			Menu:             PromptDmenu,
			ClipboardTool:    ClipboardToolXsel,
			ClipboardTimeout: 15,
			CacheTimeout:     60,
		},
		Style: ConfigurationStyle{
			PasswordBackground: "black",
			TextPassword:       "Password",
			TextMenu:           "Select",
			TextEntry:          "Entry",
			TextField:          "Field",
			FormatEntry:        "{Title} - {UserName}",
		},
		Database: ConfigurationDatabase{
			FieldOrder:      "Password UserName URL",
			FillOtherFields: true,
		},
	}
}

// LoadConfig loads the configuration into Configuration
func (c *Configuration) LoadConfig() error {
	// Get file from config folder
	viper.SetConfigType("toml")
	viper.SetConfigFile(filepath.Join(os.Getenv("HOME"), ".config/kpmenu/config"))

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Failed to read config file: %s", err)
	} else {
		// Unmarshal general
		if err := viper.UnmarshalKey("general", &c.General); err != nil {
			return NewErrorParseConfiguration("failed to parse config file (general): %v", err)
		}

		// Unmarshal style
		if err := viper.UnmarshalKey("style", &c.Style); err != nil {
			return NewErrorParseConfiguration("failed to parse config file (style): %v", err)
		}

		// Unmarshal database
		if err := viper.UnmarshalKey("database", &c.Database); err != nil {
			return NewErrorParseConfiguration("failed to parse config file (database): %v", err)
		}
	}
	return nil
}

// InitializeFlags prepare cli flags
func (c *Configuration) InitializeFlags() {
	// Flags
	flag.BoolVar(&c.Flags.Daemon, "daemon", false, "Start kpmenu directly as daemon")
	flag.BoolVarP(&c.Flags.Version, "version", "v", false, "Show kpmenu version")

	// General
	flag.StringVarP(&c.General.Menu, "menu", "m", c.General.Menu, "Choose which menu to use")
	flag.StringVar(&c.General.ClipboardTool, "clipboardTool", c.General.ClipboardTool, "Choose which clipboard tool to use")
	flag.IntVarP(&c.General.ClipboardTimeout, "clipboardTime", "c", c.General.ClipboardTimeout, "Timeout of clipboard in seconds (0 = no timeout)")
	flag.BoolVarP(&c.General.NoCache, "nocache", "n", c.General.NoCache, "Disable caching of database")
	flag.BoolVar(&c.General.CacheOneTime, "cacheOneTime", c.General.CacheOneTime, "Cache the database only the first time")
	flag.IntVar(&c.General.CacheTimeout, "cacheTimeout", c.General.CacheTimeout, "Timeout of cache in seconds")

	// Executable
	flag.StringVar(&c.Executable.CustomPromptPassword, "customPromptPassword", c.Executable.CustomPromptPassword, "Custom executable for prompt password")
	flag.StringVar(&c.Executable.CustomPromptMenu, "customPromptMenu", c.Executable.CustomPromptMenu, "Custom executable for prompt menu")
	flag.StringVar(&c.Executable.CustomPromptEntries, "customPromptEntries", c.Executable.CustomPromptEntries, "Custom executable for prompt entries")
	flag.StringVar(&c.Executable.CustomPromptFields, "customPromptFields", c.Executable.CustomPromptFields, "Custom executable for prompt fields")
	flag.StringVar(&c.Executable.CustomClipboardCopy, "customClipboardCopy", c.Executable.CustomClipboardCopy, "Custom executable for clipboard copy")
	flag.StringVar(&c.Executable.CustomClipboardPaste, "customClipboardPaste", c.Executable.CustomClipboardPaste, "Custom executable for clipboard paste")

	// Style
	flag.StringVar(&c.Style.PasswordBackground, "passwordBackground", c.Style.PasswordBackground, "Color of dmenu background and text for password selection, used to hide password typing")
	flag.StringVar(&c.Style.TextPassword, "textPassword", c.Style.TextPassword, "Label for password selection")
	flag.StringVar(&c.Style.TextMenu, "textMenu", c.Style.TextMenu, "Label for menu selection")
	flag.StringVar(&c.Style.TextEntry, "textEntry", c.Style.TextEntry, "Label for entry selection")
	flag.StringVar(&c.Style.TextField, "textField", c.Style.TextField, "Label for field selection")
	flag.StringVar(&c.Style.ArgsPassword, "argsPassword", c.Style.ArgsPassword, "Additional arguments for dmenu at password selection, separated by a space")
	flag.StringVar(&c.Style.ArgsMenu, "argsMenu", c.Style.ArgsMenu, "Additional arguments for dmenu at menu selection, separated by a space")
	flag.StringVar(&c.Style.ArgsEntry, "argsEntry", c.Style.ArgsEntry, "Additional arguments for dmenu at entry selection, separated by a space")
	flag.StringVar(&c.Style.ArgsField, "argsField", c.Style.ArgsField, "Additional arguments for dmenu at field selection, separated by a space")

	// Database
	flag.StringVarP(&c.Database.Database, "database", "d", c.Database.Database, "Path to the KeePass database")
	flag.StringVarP(&c.Database.KeyFile, "keyfile", "k", c.Database.KeyFile, "Path to the database keyfile")
	flag.StringVarP(&c.Database.Password, "password", "p", c.Database.Password, "Password of the database")
	flag.StringVar(&c.Database.FieldOrder, "fieldOrder", c.Database.FieldOrder, "String order of fields to show on field selection")
	flag.BoolVar(&c.Database.FillOtherFields, "fillOtherFields", c.Database.FillOtherFields, "Enable fill of remaining fields")
	flag.StringVar(&c.Database.FillBlacklist, "fillBlacklist", c.Database.FillBlacklist, "String of blacklisted fields that won't be shown")
}

// ParseFlags parses cli flags with given arguments
func (c *Configuration) ParseFlags(args []string) {
	flag.CommandLine.Parse(args)
}

// ErrParseConfiguration is the error return if the configuration loading fails
type ErrParseConfiguration struct {
	Message       string
	OriginalError error
}

// NewErrorParseConfiguration initializes the error
func NewErrorParseConfiguration(message string, err error) ErrParseConfiguration {
	return ErrParseConfiguration{
		Message:       message,
		OriginalError: err,
	}
}

func (err ErrParseConfiguration) Error() string {
	return fmt.Sprintf(err.Message, err.OriginalError)
}
