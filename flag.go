package cli

import (
	"fmt"
)

// Flag option of console
type Flag struct {
	Name     string   // Like 'help' equal --help
	Short    string   // Like 'h' equal -h
	Value    string   // default value
	Param    string   // Like 'path' equal --target=<path>
	Usage    string   // describe
	Required bool     // required field
	Multiple bool     // enable multiple options
	used     bool     // appear in command line
	options  []string // command line option
}

func (f *Flag) Key() string {
	return f.Name
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

// FullName reutrn name with param such as --target=<path>
func (f *Flag) FullName() string {
	if f.Name == "" {
		return ""
	}

	if f.Param != "" {
		return fmt.Sprintf("%s=<%s>", f.Name, f.Param)
	}

	return f.Name
}
