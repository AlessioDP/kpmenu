package kpmenulib

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/google/shlex"
)

// CopyToClipboard copies text into the clipboar
func CopyToClipboard(menu *Menu, text string) error {
	var cmd *exec.Cmd
	// Execute xsel/wl-clipboard/custom to update clipboard
	switch menu.Configuration.General.ClipboardTool {
	case ClipboardToolXsel:
		cmd = exec.Command("xsel", "-ib")
	case ClipboardToolWlclipboard:
		cmd = exec.Command("wl-copy")
	case ClipboardToolCustom:
		customCommand, err := shlex.Split(menu.Configuration.Executable.CustomClipboardCopy)
		if err != nil {
			return errors.New("failed to execute custom clipboard copy tool")
		}
		if len(customCommand) == 0 {
			return errors.New("the custom clipboard copy executable is empty")
		} else if len(customCommand) == 1 {
			cmd = exec.Command(customCommand[0])
		} else {
			cmd = exec.Command(customCommand[0], customCommand[1:]...)
		}
	}

	// Set sdtin
	cmd.Stdin = strings.NewReader(text)

	// Run exec
	err := cmd.Run()

	return err
}

// GetClipboard gets the current clipboard
func GetClipboard(menu *Menu) (string, error) {
	var out bytes.Buffer
	var cmd *exec.Cmd

	// Execute xsel/wl-clipboard/custom to get clipboard
	switch menu.Configuration.General.ClipboardTool {
	case ClipboardToolXsel:
		cmd = exec.Command("xsel", "-b")
	case ClipboardToolWlclipboard:
		cmd = exec.Command("wl-paste", "-n")
	case ClipboardToolCustom:
		customCommand, err := shlex.Split(menu.Configuration.Executable.CustomClipboardPaste)
		if err != nil {
			return "", errors.New("failed to execute custom executable paste tool")
		}
		if len(customCommand) == 0 {
			return "", errors.New("the custom clipboard paste executable is empty")
		} else if len(customCommand) == 1 {
			cmd = exec.Command(customCommand[0])
		} else {
			cmd = exec.Command(customCommand[0], customCommand[1:]...)
		}
	}

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
			currentClipboard, err := GetClipboard(menu)

			if err == nil {
				if text == currentClipboard {
					var cmd *exec.Cmd
					// Execute clean clipboard
					switch menu.Configuration.General.ClipboardTool {
					case ClipboardToolXsel:
						cmd = exec.Command("xsel", "-bc")
					case ClipboardToolWlclipboard:
						cmd = exec.Command("wl-copy", "-c")
					case ClipboardToolCustom:
						customCommand, err := shlex.Split(menu.Configuration.Executable.CustomClipboardClean)
						if err == nil {
							if len(customCommand) == 0 {
								cmd = nil
							} else if len(customCommand) == 1 {
								cmd = exec.Command(customCommand[0])
							} else {
								cmd = exec.Command(customCommand[0], customCommand[1:]...)
							}
						}
					}

					// Run exec
					if cmd != nil {
						if err = cmd.Run(); err != nil {
							log.Printf("failed to clean '%s' clipboard: %s", menu.Configuration.General.ClipboardTool, err)
						} else {
							log.Printf("cleaned clipboard")
						}
					}
				} else {
					log.Printf("clean clipboard cancelled because its changed")
				}
			} else {
				log.Printf("failed to get '%s' clipboard: %s", menu.Configuration.General.ClipboardTool, err)
			}
		}()
	}
}
