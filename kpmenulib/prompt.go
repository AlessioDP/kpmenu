package kpmenulib

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	
	"github.com/google/shlex"
)

// MenuSelection is an enum used for prompt menu selection
type MenuSelection int

// MenuSelections enum values
const (
	MenuShow   = MenuSelection(iota) // Show entries
	MenuReload                       // Reload database
	MenuExit                         // Exit
)

var menuSelections = [...]string{
	"Show entries",
	"Reload database",
	"Exit",
}

type entryItem struct {
	Title string
	Entry *Entry
}

// ErrorPrompt is a structure that handle an error of dmenu/rofi
type ErrorPrompt struct {
	Cancelled bool
	Error     error
}

// PromptPassword executes dmenu to ask for database password
// Returns the written password
func PromptPassword(menu *Menu) (string, ErrorPrompt) {
	// Prepare dmenu/rofi
	var command []string
	switch menu.Configuration.General.Menu {
	case "rofi":
		command = []string{
			"rofi",
			"-i",
			"-dmenu",
			"-p", menu.Configuration.Style.TextPassword,
			"-password",
		}
	case "wofi":
		command = []string{
			"wofi",
			"-i",
			"-d",
			"-p", menu.Configuration.Style.TextPassword,
			"--password",
		}
	case "dmenu":
		command = []string{
			"dmenu",
			"-i",
			"-p", menu.Configuration.Style.TextPassword,
			"-nb", menu.Configuration.Style.PasswordBackground,
			"-nf", menu.Configuration.Style.PasswordBackground,
		}
	case "custom":
		var err error
		command, err = shlex.Split(menu.Configuration.Executable.CustomPromptPassword)
		if err != nil {
			var errorPrompt ErrorPrompt
			errorPrompt.Error = fmt.Errorf("failed to parse custom prompt password, exiting")
			return "", errorPrompt
		}
	}

	// Add custom arguments
	if menu.Configuration.Style.ArgsPassword != "" {
		command = append(command, strings.Split(menu.Configuration.Style.ArgsPassword, " ")...)
	}

	// Execute prompt
	return executePrompt(command, nil)
}

// PromptMenu executes dmenu to ask for menu selection
// Returns the MenuSelection chosen
func PromptMenu(menu *Menu) (MenuSelection, ErrorPrompt) {
	var selection MenuSelection
	var input strings.Builder

	// Prepare dmenu/rofi
	var command []string
	switch menu.Configuration.General.Menu {
	case "rofi":
		command = []string{
			"rofi",
			"-i",
			"-dmenu",
			"-p", menu.Configuration.Style.TextMenu,
		}
	case "wofi":
		command = []string{
			"wofi",
			"-i",
			"-d",
			"-p", menu.Configuration.Style.TextMenu,
		}
	case "dmenu":
		command = []string{
			"dmenu",
			"-i",
			"-p", menu.Configuration.Style.TextMenu,
		}
	case "custom":
		var err error
		command, err = shlex.Split(menu.Configuration.Executable.CustomPromptMenu)
		if err != nil {
			var errorPrompt ErrorPrompt
			errorPrompt.Cancelled = true
			errorPrompt.Error = fmt.Errorf("failed to parse custom prompt menu, exiting")
			return 0, errorPrompt
		}
	}

	// Add custom arguments
	if menu.Configuration.Style.ArgsMenu != "" {
		command = append(command, strings.Split(menu.Configuration.Style.ArgsMenu, " ")...)
	}

	// Prepare input (dmenu items)
	for _, e := range menuSelections {
		input.WriteString(e + "\n")
	}

	// Execute prompt
	result, err := executePrompt(command, strings.NewReader(input.String()))
	if err.Error == nil && !err.Cancelled {
		// Get selected menu item
		for ind, sel := range menuSelections {
			// Match for entry title and selected entry
			if sel == result {
				selection = MenuSelection(ind)
				break
			}
		}
	}
	return selection, err
}

// PromptEntries executes dmenu to ask for an entry selection
// Returns the selected entry
func PromptEntries(menu *Menu) (*Entry, ErrorPrompt) {
	var entry Entry
	var input strings.Builder

	// Prepare dmenu/rofi
	var command []string
	switch menu.Configuration.General.Menu {
	case "rofi":
		command = []string{
			"rofi",
			"-i",
			"-dmenu",
			"-p", menu.Configuration.Style.TextEntry,
		}
	case "wofi":
		command = []string{
			"wofi",
			"-i",
			"-d",
			"-p", menu.Configuration.Style.TextEntry,
		}
	case "dmenu":
		command = []string{
			"dmenu",
			"-i",
			"-p", menu.Configuration.Style.TextEntry,
		}
	case "custom":
		var err error
		command, err = shlex.Split(menu.Configuration.Executable.CustomPromptEntries)
		if err != nil {
			var errorPrompt ErrorPrompt
			errorPrompt.Error = fmt.Errorf("failed to parse custom prompt entries, exiting")
			return nil, errorPrompt
		}
	}

	// Add custom arguments
	if menu.Configuration.Style.ArgsEntry != "" {
		command = append(command, strings.Split(menu.Configuration.Style.ArgsEntry, " ")...)
	}

	// Prepare a list of entries
	// Identified by the formatted title and the entry pointer
	var listEntries []entryItem
	reg, err := regexp.Compile(`{[a-zA-Z]+\}`)
	if err != nil {
		return &entry, ErrorPrompt{
			Cancelled: false,
			Error:     err,
		}
	}
	for i, e := range menu.Database.Entries {
		// Format entry
		title := menu.Configuration.Style.FormatEntry
		matches := reg.FindAllString(title, -1)

		// Replace every match
		for _, match := range matches {
			valueType := match[1 : len(match)-1] // Removes { and }
			value := ""                          // By default empty value
			vd := e.FullEntry.GetContent(valueType)
			if vd != "" {
				value = vd
			}
			title = strings.Replace(title, match, value, -1)
		}
		// Be sure to point on the right entry, do not point to the local e
		listEntries = append(listEntries, entryItem{Title: title, Entry: &menu.Database.Entries[i]})
	}

	// Prepare input (dmenu items)
	for _, e := range listEntries {
		input.WriteString(e.Title + "\n")
	}

	// Execute prompt
	result, errPrompt := executePrompt(command, strings.NewReader(input.String()))
	if errPrompt.Error == nil && !errPrompt.Cancelled {
		// Get selected entry
		for _, e := range listEntries {
			if e.Title == result {
				entry = *e.Entry
				break
			}
		}
	}
	return &entry, errPrompt
}

// PromptFields executes dmenu to ask for a field selection
// Returns the selected field value as string
func PromptFields(menu *Menu, entry *Entry) (string, ErrorPrompt) {
	var value string
	var input strings.Builder

	// Prepare dmenu/rofi
	var command []string
	switch menu.Configuration.General.Menu {
	case "rofi":
		command = []string{
			"rofi",
			"-i",
			"-dmenu",
			"-p", menu.Configuration.Style.TextEntry,
		}
	case "wofi":
		command = []string{
			"wofi",
			"-i",
			"-d",
			"-p", menu.Configuration.Style.TextEntry,
		}
	case "dmenu":
		command = []string{
			"dmenu",
			"-i",
			"-p", menu.Configuration.Style.TextField,
		}
	case "custom":
		var err error
		command, err = shlex.Split(menu.Configuration.Executable.CustomPromptFields)
		if err != nil {
			var errorPrompt ErrorPrompt
			errorPrompt.Error = fmt.Errorf("failed to parse custom prompt fields, exiting")
			return "", errorPrompt
		}
	}

	// Add custom arguments
	if menu.Configuration.Style.ArgsField != "" {
		command = append(command, strings.Split(menu.Configuration.Style.ArgsField, " ")...)
	}

	fields := []string{}
	fieldsOrder := strings.Split(menu.Configuration.Database.FieldOrder, " ")

	// Populate ordered fields
	for _, f := range fieldsOrder {
		if entry.FullEntry.GetContent(f) != "" {
			fields = append(fields, f)
		}
	}

	// Populate with filling fields
	if menu.Configuration.Database.FillOtherFields {
		blacklistFields := strings.Split(menu.Configuration.Database.FillBlacklist, " ")

		for _, v := range entry.FullEntry.Values {
			if !contains(fields, v.Key) && !contains(blacklistFields, v.Key) {
				if v.Value.Content != "" {
					fields = append(fields, v.Key)
				}
			}
		}
	}

	// Prepare input (dmenu items)
	for _, f := range fields {
		input.WriteString(f + "\n")
	}

	// Execute prompt
	result, err := executePrompt(command, strings.NewReader(input.String()))
	if err.Error == nil && !err.Cancelled {
		// Check that the result is valid
		if contains(fields, result) {
			// Get field value
			for _, v := range entry.FullEntry.Values {
				if result == v.Key {
					value = v.Value.Content
					break
				}
			}
		}
	}
	return value, err
}

func executePrompt(command []string, input *strings.Reader) (result string, errorPrompt ErrorPrompt) {
	var out bytes.Buffer
	var outErr bytes.Buffer
	var cmd *exec.Cmd
	if len(command) == 0 {
		errorPrompt.Cancelled = true
		errorPrompt.Error = fmt.Errorf("the custom prompt command is empty")
		return
	} else if len(command) == 1 {
		cmd = exec.Command(command[0])
	} else {
		cmd = exec.Command(command[0], command[1:]...)
	}
	

	// Set stdout to out var
	cmd.Stdout = &out
	cmd.Stderr = &outErr
	if input != nil {
		cmd.Stdin = input
	}

	// Run exec
	err := cmd.Run()
	if err != nil {
		if outErr.String() != "" {
			errorPrompt.Error = fmt.Errorf(
				"the command %s returned %s: %s",
				command,
				err,
				outErr.String(),
			)
		} else {
			errorPrompt.Cancelled = true
		}
	}
	// Trim output
	result = strings.TrimRight(out.String(), "\n")
	return
}

func contains(array []string, value string) bool {
	for _, n := range array {
		if value == n {
			return true
		}
	}
	return false
}

func (el MenuSelection) String() string {
	return menuSelections[el]
}
