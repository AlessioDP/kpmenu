package kpmenulib

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
	"time"
)

// CopyToClipboard copies text into the clipboar
func CopyToClipboard(text string) error {
	// Execute xsel to update clipboard
	cmd := exec.Command("xsel", "-ib")

	// Set sdtin
	cmd.Stdin = strings.NewReader(text)

	// Run exec
	err := cmd.Run()

	return err
}

// GetClipboard gets the current clipboard
func GetClipboard() (string, error) {
	var out bytes.Buffer

	// Execute xsel to get clipboard
	cmd := exec.Command("xsel", "-b")

	// Set stdout
	cmd.Stdout = &out

	// Run exec
	err := cmd.Run()

	return out.String(), err
}

// CleanClipboard cleans the clipboard, if not changed
func CleanClipboard(menu *Menu, text string) {
	if menu.Configuration.General.ClipboardTimeout > 0 {
		// Goroutine
		// Its async so any error will be printed
		menu.WaitGroup.Add(1)
		go func() {
			defer menu.WaitGroup.Done()
			// Sleep for X seconds
			time.Sleep(time.Duration(menu.Configuration.General.ClipboardTimeout) * time.Second)

			// Execute GetClipboard to match old and current cliboard
			// Clean clipboard only if contains the field value
			currentClipboard, err := GetClipboard()
			if err == nil {
				if text == currentClipboard {
					// Execute xsel to clean clipboard
					cmd := exec.Command("xsel", "-bc")

					// Run exec
					if err = cmd.Run(); err != nil {
						log.Printf("failed to clean clipboard: %s", err)
					} else {
						log.Printf("cleaned clipboard")
					}
				} else {
					log.Printf("clean clipboard cancelled because its changed")
				}
			} else {
				log.Printf("failed to get the clipboard: %s", err)
			}
		}()
	}
}
