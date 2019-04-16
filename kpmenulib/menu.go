package kpmenulib

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// Menu is the main structure of kpmenu
type Menu struct {
	CacheStart    time.Time       // Cache start time
	CliArguments  []string        // Arguments of kpmenu
	Configuration *Configuration  // Configuration of kpmenu
	Database      *Database       // Database
	WaitGroup     *sync.WaitGroup // WaitGroup used for goroutines
}

// NewMenu initializes a Menu struct
func NewMenu() *Menu {
	return &Menu{
		CliArguments:  os.Args[1:],
		Configuration: NewConfiguration(),
		Database:      NewDatabase(),
		WaitGroup:     new(sync.WaitGroup),
	}
}

// OpenDatabase asks for password and populates the database
func (m *Menu) OpenDatabase() *ErrorDatabase {
	// Check if there is already a password/key set
	if len(m.Database.Keepass.Credentials.Passphrase) == 0 &&
		len(m.Database.Keepass.Credentials.Key) == 0 {
		// Get password from config otherwise ask for it
		password := m.Configuration.Database.Password
		if password == "" {
			// Get password from user
			pw, err := PromptPassword(m)
			if !err.Cancelled {
				if err.Error != nil {
					return NewErrorDatabase("failed to get password from dmenu: %s", err.Error, true)
				}
			} else {
				// Exit because cancelled
				return NewErrorDatabase("exiting because user cancelled password prompt", nil, true)
			}
			password = pw
		}

		// Add credentials into the database
		m.Database.AddCredentialsToDatabase(m.Configuration, password)
	}

	// Open database
	if err := m.Database.OpenDatabase(m.Configuration); err != nil {
		return NewErrorDatabase("failed to open database: %s", err, true)
	}

	// Get entries of database
	m.Database.IterateDatabase()

	// Set database as loaded
	m.Database.Loaded = true

	return nil
}

// OpenMenu executes dmenu to interface the user with the database
func (m *Menu) OpenMenu() *ErrorDatabase {
	// Prompt for menu selection
	selectedMenu, err := PromptMenu(m)
	if err.Cancelled {
		if err.Error != nil {
			return NewErrorDatabase("failed to select menu item: %s", err.Error, false)
		}
		// Cancelled
		return NewErrorDatabase("", nil, false)
	}
	switch selectedMenu {
	case MenuShow:
		return m.entrySelection()
	case MenuReload:
		log.Printf("reloading database")
		if err := m.OpenDatabase(); err != nil {
			return err
		}
		return m.OpenMenu()
	case MenuExit:
		return NewErrorDatabase("exiting", nil, true)
	}
	return nil
}

func (m *Menu) entrySelection() *ErrorDatabase {
	// Prompt for entry selection
	selectedEntry, err := PromptEntries(m)
	if err.Cancelled {
		if err.Error != nil {
			return NewErrorDatabase("failed to select entry: %s", err.Error, false)
		}
		// Cancelled
		return NewErrorDatabase("", nil, false)
	}
	if selectedEntry == nil {
		// Entry not found
		return NewErrorDatabase("selected entry not found", nil, false)
	}

	// Prompt for field selection
	fieldValue, err := PromptFields(m, selectedEntry)
	if err.Cancelled {
		if err.Error != nil {
			return NewErrorDatabase("failed to select field: %s", err.Error, false)
		}
		// Cancelled
		return NewErrorDatabase("", nil, false)
	}
	if fieldValue == "" {
		// Field not found
		return NewErrorDatabase("selected field not found", nil, false)
	}

	// Copy to clipboard
	if err := CopyToClipboard(m, fieldValue); err != nil {
		return NewErrorDatabase("failed to use clipboard manager to update clipboard: %s", err, true)
	}
	log.Printf("copied field into the clipboard")

	// Clean clipboard (goroutine)
	CleanClipboard(m, fieldValue)
	return nil
}

// ErrorDatabase is an error that can be fatal or non-fatal
type ErrorDatabase struct {
	Message       string
	OriginalError error
	Fatal         bool
}

// NewErrorDatabase makes an ErrorDatabase
func NewErrorDatabase(message string, err error, fatal bool) *ErrorDatabase {
	return &ErrorDatabase{
		Message:       message,
		OriginalError: err,
		Fatal:         fatal,
	}
}

func (err *ErrorDatabase) String() string {
	if err.OriginalError != nil {
		return fmt.Sprintf(err.Message, err.OriginalError)
	}
	return err.Message
}
