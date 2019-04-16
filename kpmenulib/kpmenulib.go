package kpmenulib

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"
)

// Version is the version of kpmenu
const Version = "1.2.0-SNAPSHOT"

// Initialize is the function that initialize a menu, handle config and parse cli arguments
func Initialize() *Menu {
	// Initialize menu
	menu := NewMenu()

	// Load configuration
	if err := handleConfiguration(menu, false); err != nil {
		log.Fatal(err)
		return nil
	}

	// Load cli flags
	if quit := handleCli(menu); quit == true {
		return nil
	}

	// Set start cache time, if not a daemon
	if !menu.Configuration.Flags.Daemon && !menu.Configuration.General.NoCache {
		menu.CacheStart = time.Now()
	}

	return menu
}

// Execute is the function used to open the database (if necessary) and open the menu
// returns true if the program should exit
func Execute(menu *Menu) bool {
	// Open database
	if menu.Database.Loaded == false {
		if err := menu.OpenDatabase(); err != nil {
			log.Print(err)
			return err.Fatal
		}
	}

	// Open menu
	if err := menu.OpenMenu(); err != nil {
		log.Print(err)
		return err.Fatal
	}

	// Non-fatal exit
	return false
}

// Show checks if the database configuration is changed, if so it will re-open the database
// returns true if the program should exit
func Show(menu *Menu) bool {
	// Be sure that the database configuration is the same, otherwise a Run is necessary
	copiedDatabase := menu.Configuration.Database

	// Re handle configuration and update it if changed
	if err := handleConfiguration(menu, true); err != nil {
		log.Print(err)
		return true
	}
	menu.Configuration.ParseFlags(menu.CliArguments)

	// If something related to the database is changed we must re-open it, or exit true
	if copiedDatabase.Database != menu.Configuration.Database.Database ||
		copiedDatabase.KeyFile != menu.Configuration.Database.KeyFile ||
		copiedDatabase.Password != menu.Configuration.Database.Password {
		menu.Database.Loaded = false
		log.Printf("database configuration is changed, re-opening the database")
	}

	// Check if the cache is not timed out, if not a daemon
	if !menu.Configuration.Flags.Daemon {
		if menu.Configuration.General.NoCache {
			// Cache disabled
			menu.Database.Loaded = false
			log.Printf("no cache flag is set, re-opening the database")
		} else if (menu.CacheStart == time.Time{}) {
			// Cache enabled via client call
			menu.Database.Loaded = false
			log.Printf("cache start time not set, re-opening the database")
		} else {
			// Cache exists
			difference := int(time.Now().Sub(menu.CacheStart).Seconds())
			if difference < menu.Configuration.General.CacheTimeout {
				// Cache is valid
				if !menu.Configuration.General.CacheOneTime {
					// Set new cache start if cache one time is false
					menu.CacheStart = time.Now()
				}
			} else {
				// Cache timed out
				menu.Database.Loaded = false
				log.Printf("cache timed out, re-opening the database")
			}
		}
	}

	return Execute(menu)
}

func handleConfiguration(menu *Menu, parseOnly bool) error {
	// Load configuration
	if err := menu.Configuration.LoadConfig(); err != nil {
		return err
	}

	// Initialize flags if parseOnly false
	if !parseOnly {
		menu.Configuration.InitializeFlags()
	}
	menu.Configuration.ParseFlags(menu.CliArguments)

	if err := checkFlags(menu); err != nil {
		return err
	}
	return nil
}

func checkFlags(menu *Menu) error {
	// Check if database has been selected
	if menu.Configuration.Database.Database == "" {
		// Database not found
		return errors.New("you must select a database with -d or via config")
	}

	// Check if rofi is installed
	if menu.Configuration.General.UseRofi {
		cmd := exec.Command("which", "rofi")
		err := cmd.Run()
		if err != nil {
			log.Printf("rofi not found, using dmenu")
			menu.Configuration.General.UseRofi = false
		}
	}

	// Check if dmenu is installed
	if !menu.Configuration.General.UseRofi {
		cmd := exec.Command("which", "dmenu")
		err := cmd.Run()
		if err != nil {
			return errors.New("dmenu not found, exiting")
		}
	}

	if menu.Configuration.General.ClipboardTool == ClipboardToolWlclipboard {
		// Check if wl-clipboard is installed
		cmd := exec.Command("which", "wl-copy")
		err := cmd.Run()
		if err != nil {
			return errors.New("wl-clipboard not found, exiting")
		}
	} else {
		// Check if xsel is installed
		cmd := exec.Command("which", "xsel")
		err := cmd.Run()
		if err != nil {
			return errors.New("xsel not found, exiting")
		}
	}
	return nil
}

// Returns true if the program should exit
func handleCli(menu *Menu) bool {
	// Check for version flag
	if menu.Configuration.Flags.Version {
		// Print version and exit
		fmt.Println(Version)
		return true
	}
	return false
}
