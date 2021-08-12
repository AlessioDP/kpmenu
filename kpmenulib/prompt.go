package kpmenulib

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

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

	var hasOTP bool
	// Populate with filling fields
	if menu.Configuration.Database.FillOtherFields {
		blacklistFields := strings.Split(menu.Configuration.Database.FillBlacklist, " ")

		for _, v := range entry.FullEntry.Values {
			if !menu.Configuration.General.NoOTP && (v.Key == OTP || v.Key == TOTPSEED) {
				hasOTP = true
				continue
			}
			if !contains(fields, v.Key) && !contains(blacklistFields, v.Key) {
				if v.Value.Content != "" {
					fields = append(fields, v.Key)
				}
			}
		}
	}

	// Prepare input (dmenu items)
	const GenerateOTP = "Generate OTP"
	for _, f := range fields {
		input.WriteString(f + "\n")
	}
	if hasOTP {
		input.WriteString(GenerateOTP + "\n")
	}

	// Execute prompt
	result, err := executePrompt(command, strings.NewReader(input.String()))
	if err.Error == nil && !err.Cancelled {
		if result == GenerateOTP {
			var ev error
			value, ev = CreateOTP(entry.FullEntry, time.Now().Unix())
			if ev != nil {
				err.Cancelled = true
				err.Error = fmt.Errorf("failed to create otp: %s", ev)
				return value, err
			}
			return value, err
		}
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

// PromptAutotype executes an external application to select an entry and then
// runs an autotype program with the entry's data.
//
// Field data is sent to the autotype child process on STDIN as TSV data. The
// first row is the key sequence.
//
//    {USERNAME}{TAB}{PASSWORD}{ENTER}
//    key <TAB> value <CR>
//
// If KeepassXC attributes for autotype exist for the record, they're used. If
// they do not exist, username & password are used, and if OTP data exists for
// the record, that is sent as the third field.
//
// STDIN is closed when all fields have been sent.
//
// Note that it is currently not possible to get the key sequence from the DB,
// so it is always the default.
func PromptAutotype(menu *Menu) ErrorPrompt {
	// Prepare autotype command
	var command []string
	command = strings.Split(menu.Configuration.Executable.CustomAutotypeWindowID, " ")

	activeWindow, errPrompt := executePrompt(command, nil)
	if errPrompt.Error != nil || errPrompt.Cancelled {
		return errPrompt
	}

	// Only a pointer to allow the nil test for a positive find
	var entry *Entry
	for _, e := range menu.Database.Entries {
		ms := e.FullEntry.GetContent("Title")
		at := e.FullEntry.AutoType.Association
		// Check if there's a window association, and if so, make sure it's a regexp
		if at != nil && at.Window != "" {
			// For regexp, remove the wrapping // and replace all star globs with .*
			if !strings.HasPrefix(at.Window, "//") {
				ms = strings.ReplaceAll(at.Window, "*", ".*")
			} else {
				if len(ms) > 2 {
					ms = ms[2:]
				}
				if len(ms) > 2 {
					ms = ms[:len(ms)-2]
				}
			}
		}
		reg, err := regexp.Compile(ms)
		if err != nil {
			continue
		}
		if reg.Match([]byte(activeWindow)) {
			entry = &e
			break
		}
	}

	if entry == nil {
		errPrompt.Error = fmt.Errorf("no autotype window match for %s", activeWindow)
		return errPrompt
	}

	var input strings.Builder
	fe := entry.FullEntry
	if fe.AutoType.Association != nil && fe.AutoType.Association.KeystrokeSequence != "" {
		input.WriteString(fe.AutoType.Association.KeystrokeSequence)
		input.WriteString("\n")
	} else if !menu.Configuration.General.NoOTP && (fe.GetContent(OTP) != "" || fe.GetContent(TOTPSEED) != "") {
		input.WriteString("{USERNAME}{TAB}{PASSWORD}{TAB}{TOTP}{ENTER}\n")
	} else {
		input.WriteString("{USERNAME}{TAB}{PASSWORD}{ENTER}\n")
	}
	input.WriteString("UserName")
	input.WriteString("\t")
	input.WriteString(entry.FullEntry.GetContent("UserName"))
	input.WriteString("\n")

	input.WriteString("Password")
	input.WriteString("\t")
	input.WriteString(entry.FullEntry.GetContent("Password"))
	input.WriteString("\n")

	if !menu.Configuration.General.NoOTP && (fe.GetContent(OTP) != "" || fe.GetContent(TOTPSEED) != "") {
		value, err := CreateOTP(fe, time.Now().Unix())
		if err != nil {
			errPrompt.Cancelled = true
			errPrompt.Error = fmt.Errorf("failed to create otp: %s", err)
			return errPrompt
		}
		input.WriteString("OTP")
		input.WriteString("\t")
		input.WriteString(value)
		input.WriteString("\n")
	}

	command = strings.Split(menu.Configuration.Executable.CustomAutotypeTyper, " ")
	_, errPrompt = executePrompt(command, strings.NewReader(input.String()))

	return errPrompt
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
