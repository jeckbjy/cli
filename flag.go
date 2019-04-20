package cli

import (
	"fmt"
	"strings"
)

// Flag arg
type Flag struct {
	key      string      // the first name
	Name     string      // split by comma like "help,h"
	Value    string      // default value
	Hint     interface{} // hint type
	Required bool        // required field
	Multiple bool        // enable multiple options
	used     bool        // appear in command line
	options  []string    // command line option
}

func (f *Flag) Key() string {
	if f.key == "" {
		index := strings.Index(f.Name, ",")
		if index != -1 {
			f.key = strings.TrimSpace(f.Name[:index])
		} else {
			f.key = f.Name
		}
	}

	return f.key
}

func (f *Flag) addOption(opt string) error {
	if f.used && !f.Multiple {
		return fmt.Errorf("option cannot multipe:%+v", f.Name)
	}

	f.used = true

	if opt != "" {
		f.options = append(f.options, opt)
	}

	return nil
}

func (f *Flag) validate() error {
	if f.Required && !f.used {
		return fmt.Errorf("option is required:%+v", f.Name)
	}

	// set default to options
	if len(f.options) == 0 && f.Value != "" {
		f.options = append(f.options, f.Value)
	}

	return nil
}

// Get return option or default value
func (f *Flag) Get() string {
	if len(f.options) > 0 {
		return f.options[0]
	}

	return ""
}

// GetAt return option or default value by index if optoin is array
func (f *Flag) GetAt(i int) string {
	if i < len(f.options) {
		return f.options[i]
	}

	return ""
}

// GetList return the options list
func (f *Flag) GetList() []string {
	return f.options
}

// Len return the length of the option
func (f *Flag) Len() int {
	return len(f.options)
}
