package kpmenulib

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// Field names, e.g. {USERNAME}
	FIELD = iota
	// Keywords, e.g. {TAB} or {ENTER}
	KEYWORD
	// Commands, e.g. {DELAY 5}
	COMMAND
	// Raw text, e.g. text not enclosed in {}
	RAW
	// Special alias characters, e.g. a ^ not in {} means "control"
	SPECIAL
)

// SeqEntry is a single token in a key sequence, denoted by the token, the
// parsed type, and any args if it is a command.
type SeqEntry struct {
	// token is the processed text, stripped of {}
	Token string
	// args is only set for COMMANDs, and will be nil otherwise
	Args []string
	// The type of the sequence entry, e.g. KEYWORD, COMMAND, etc.
	Type int
}

type SeqEntries []SeqEntry

// Sequence is a parsed sequence of tokens
type Sequence struct {
	SeqEntries
	Keylag int
}

// NewSequence returns a new Sequence instance bound to a typer. Unless mocking, this
// should normally be:
// ```
// s := NewSequence(Robot{})
// ```
func NewSequence() Sequence {
	return Sequence{
		make(SeqEntries, 0),
		50,
	}
}

// Parse processes a [Keepass autotype sequence](https://keepass.info/help/base/autotype.html)
// and returns the parsed keys in the order in which they occurred.
func (rv *Sequence) Parse(keySeq string) error {
	if _atKeySeqRE == nil {
		initKeySeqParser()
	}
	if len(keySeq) == 0 {
		return fmt.Errorf("received empty sequence")
	}
	matches := _atKeySeqRE.FindAllString(keySeq, -1)
	if len(matches) == 0 {
		return fmt.Errorf("received malformed sequence %#v", keySeq)
	}
	for _, match := range matches {
		var s SeqEntry
		switch {
		case match == "{{}":
			s.Token = "{"
			s.Type = KEYWORD
		case match == "{}}":
			s.Token = "}"
			s.Type = KEYWORD
		case match[0] == '{':
			match = strings.Trim(match, "{}")
			match = strings.Trim(match, " ")
			match = strings.Replace(match, "=", " ", 1)
			parts := strings.Split(match, " ")
			s.Token = parts[0]
			if contains(_commands, s.Token) {
				s.Type = COMMAND
				if len(parts) > 1 {
					s.Args = parts[1:]
				} else {
					s.Args = []string{}
				}
			} else if contains(_atKeywords, match) {
				s.Token = match
				s.Type = KEYWORD
			} else if len(match) == 0 {
				return fmt.Errorf("invalid key sequence {}")
			} else {
				s.Token = match
				s.Type = FIELD
			}
		default:
			s.Token = match
			s.Type = RAW
			if contains(_specials, match) {
				s.Type = SPECIAL
			}
		}
		rv.SeqEntries = append(rv.SeqEntries, s)
	}
	return nil
}

const (
	// Argument-laden keywords: DELAY, VKEY, APPACTIVATE, BEEP
	AT_KW_DELAY  = "DELAY"
	AT_KW_VKEY   = "VKEY"
	AT_KW_APPACT = "APPACTIVATE"
	AT_KW_BEEP   = "BEEP"
	// No-argument keywords
	AT_KW_CLEAR             = "CLEARFIELD"
	AT_KW_PLUS              = "PLUS"
	AT_KW_PERCENT           = "PERCENT"
	AT_KW_CARET             = "CARET"
	AT_KW_TILDE             = "TILDE"
	AT_KW_LEFTPAREN         = "LEFTPAREN"
	AT_KW_RIGHTPAREN        = "RIGHTPAREN"
	AT_KW_LEFTBRACE         = "LEFTBRACE"
	AT_KW_RIGHTBRACE        = "RIGHTBRACE"
	AT_KW_AT                = "AT"
	AT_KW_TAB               = "TAB"
	AT_KW_ENTER             = "ENTER"
	AT_KW_ARROW_UP          = "UP"
	AT_KW_ARROW_DOWN        = "DOWN"
	AT_KW_ARROW_LEFT        = "LEFT"
	AT_KW_ARROW_RIGHT       = "RIGHT"
	AT_KW_INSERT            = "INSERT"
	AT_KW_INSERT2           = "INS"
	AT_KW_DELETE            = "DELETE"
	AT_KW_DELETE2           = "DEL"
	AT_KW_HOME              = "HOME"
	AT_KW_END               = "END"
	AT_KW_PAGE_UP           = "PGUP"
	AT_KW_PAGE_DOWN         = "PGDN"
	AT_KW_SPACE             = "SPACE"
	AT_KW_BREAK             = "BREAK"
	AT_KW_CAPS_LOCK         = "CAPSLOCK"
	AT_KW_ESCAPE            = "ESC"
	AT_KW_BACKSPACE         = "BACKSPACE"
	AT_KW_KW_BACKSPACE2     = "BS"
	AT_KW_KW_BACKSPACE3     = "BKSP"
	AT_KW_WINDOWS_KEY_LEFT  = "WIN"
	AT_KW_WINDOWS_KEY_LEFT2 = "LWIN"
	AT_KW_WINDOWS_KEY_RIGHT = "RWIN"
	AT_KW_HELP              = "HELP"
	AT_KW_NUMLOCK           = "NUMLOCK"
	AT_KW_PRINTSCREEN       = "PRTSC"
	AT_KW_SCROLLLOCK        = "SCROLLLOCK"
	AT_KW_NUMPAD_PLUS       = "ADD"
	AT_KW_NUMPAD_MINUS      = "SUBTRACT"
	AT_KW_NUMPAD_MULT       = "MULTIPLY"
	AT_KW_NUMPAD_DIV        = "DIVIDE"
	AT_KW_APPS_MENU         = "APPS"
	//Numeric pad 0 to 9	{NUMPAD0} to {NUMPAD9}
	// F1 - F16	{F1} - {F16}
	// Special characters within {} -- thanks, keepass2!
	AT_SCHAR = "+%^~{}[]()"
	// Special characters *outside* of {}
	AT_CH_SHIFT      = "+"
	AT_CH_CTRL       = "^"
	AT_CH_ALT        = "%"
	AT_CH_WINDOWSKEY = "@"
	AT_CH_ENTER2     = "~"
)

var _commands []string = []string{AT_KW_DELAY, AT_KW_VKEY, AT_KW_APPACT, AT_KW_BEEP}

var _specials []string = []string{AT_CH_SHIFT, AT_CH_CTRL, AT_CH_ALT, AT_CH_WINDOWSKEY, AT_CH_ENTER2}

// This has been benchmarked. An O(M*N) array search is faster than either map or regex.
var _atKeywords []string = []string{
	AT_KW_CLEAR, AT_KW_PLUS, AT_KW_PERCENT,
	AT_KW_CARET, AT_KW_TILDE, AT_KW_LEFTPAREN, AT_KW_RIGHTPAREN, AT_KW_LEFTBRACE, AT_KW_RIGHTBRACE,
	AT_KW_AT, AT_KW_TAB, AT_KW_ENTER, AT_KW_ARROW_UP, AT_KW_ARROW_DOWN, AT_KW_ARROW_LEFT,
	AT_KW_ARROW_RIGHT, AT_KW_INSERT, AT_KW_INSERT2, AT_KW_DELETE, AT_KW_DELETE2, AT_KW_HOME,
	AT_KW_END, AT_KW_PAGE_UP, AT_KW_PAGE_DOWN, AT_KW_SPACE, AT_KW_BREAK, AT_KW_CAPS_LOCK,
	AT_KW_ESCAPE, AT_KW_BACKSPACE, AT_KW_KW_BACKSPACE2, AT_KW_KW_BACKSPACE3, AT_KW_WINDOWS_KEY_LEFT,
	AT_KW_WINDOWS_KEY_LEFT2, AT_KW_WINDOWS_KEY_RIGHT, AT_KW_HELP, AT_KW_NUMLOCK, AT_KW_PRINTSCREEN,
	AT_KW_SCROLLLOCK, AT_KW_NUMPAD_PLUS, AT_KW_NUMPAD_MINUS, AT_KW_NUMPAD_MULT, AT_KW_NUMPAD_DIV,
	AT_KW_APPS_MENU,
}

// Breaks a sequence string into tokens
var _atKeySeqRE *regexp.Regexp

// Don't make this init() -- it will (unnecessarily) slow down client calls. You probably
// don't want to call it; it's called by `Parse()` if necessary
func initKeySeqParser() {
	_atKeySeqRE = regexp.MustCompile("\\{[^\\}]+\\}|[^\\{\\}]+|\\s+|\\{\\{\\}|\\{\\}\\}")

	// Add special characters in {}
	for _, c := range strings.Split(AT_SCHAR, "") {
		_atKeywords = append(_atKeywords, c)
	}
	// NUMPAD0 - NUMPAD9
	for i := 0; i < 10; i++ {
		_atKeywords = append(_atKeywords, fmt.Sprintf("NUMPAD%d", i))
	}
	// F1 - F16
	for i := 1; i < 17; i++ {
		_atKeywords = append(_atKeywords, fmt.Sprintf("F%d", i))
	}
}

// contains answers: does an array of strings contain a string?
func contains(kws []string, s string) bool {
	for _, keyword := range kws {
		if keyword == s {
			return true
		}
	}
	return false
}
