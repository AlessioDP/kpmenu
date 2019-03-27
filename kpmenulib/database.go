package kpmenulib

import (
	"log"
	"os"

	"github.com/tobischo/gokeepasslib/v2"
)

// Database contains the KeePass database and its entry list
type Database struct {
	Loaded  bool
	Keepass *gokeepasslib.Database
	Entries []Entry
}

// Entry is a container for keepass entry
type Entry struct {
	UUID      gokeepasslib.UUID
	FullEntry gokeepasslib.Entry
}

// NewDatabase initializes the Database struct
func NewDatabase() *Database {
	return &Database{
		Loaded:  false,
		Keepass: gokeepasslib.NewDatabase(),
	}
}

// AddCredentialsToDatabase adds credentials into gokeepasslib credentials struct
func (db *Database) AddCredentialsToDatabase(cfg *Configuration, password string) {
	// Get credentials
	if password != "" && cfg.Database.KeyFile != "" {
		// Both password & keyfile
		db.Keepass.Credentials, _ = gokeepasslib.NewPasswordAndKeyCredentials(password, cfg.Database.KeyFile)
		log.Printf("credentials: password + keyfile")
	} else if password != "" {
		// Only password
		db.Keepass.Credentials = gokeepasslib.NewPasswordCredentials(password)
		log.Printf("credentials: password")
	} else if cfg.Database.KeyFile != "" {
		// Only keyfile
		db.Keepass.Credentials, _ = gokeepasslib.NewKeyCredentials(cfg.Database.KeyFile)
		log.Printf("credentials: keyfile")
	}
}

// OpenDatabase decodes the database with the given configuration
func (db *Database) OpenDatabase(cfg *Configuration) error {
	// Open database file
	file, err := os.Open(cfg.Database.Database)
	if err == nil {
		err = gokeepasslib.NewDecoder(file).Decode(db.Keepass)
		if err == nil {
			err = db.Keepass.UnlockProtectedEntries()
		}
	}
	return err
}

// IterateDatabase iterates the database and makes a list of entries
func (db *Database) IterateDatabase() {
	var entries []Entry
	for _, sub := range db.Keepass.Content.Root.Groups {
		entries = append(entries, iterateGroup(sub)...)
	}
	db.Entries = entries
}

func iterateGroup(kpGroup gokeepasslib.Group) []Entry {
	var entries []Entry
	// Get entries of the current group
	for _, kpEntry := range kpGroup.Entries {
		// Insert entry
		entries = append(entries, Entry{
			UUID:      kpEntry.UUID,
			FullEntry: kpEntry,
		})
		//(*entries)[uuid] = Entry{FullEntry: kpEntry}
	}

	// Continue to iterate subgroups
	for _, sub := range kpGroup.Groups {
		entries = append(entries, iterateGroup(sub)...)
	}
	return entries
}
