package cli

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	styleUnknown = 0 // bad
	styleWindow  = 1 // '/'
	styleSingle  = 2 // '-'
	styleDouble  = 3 // '--'
)

const (
	AppCommandName = ""
)

// ProcessBar
// http://www.gnu.org/software/bash/manual/bash.html#Programmable-Completion
// https://github.com/cheggaaa/pb
// https://github.com/vbauerster/mpb
// https://github.com/schollz/progressbar

// prompt
// https://github.com/c-bata/go-prompt
// https://github.com/manifoldco/promptui
// https://github.com/justjanne/powerline-go

// cli
// https://github.com/urfave/cli
// https://github.com/spf13/cobra

// flags
// https://github.com/jessevdk/go-flags

// App build a git style cli
type App struct {
	Name     string
	Flags    []*Flag           // Top Flag
	Commands []*Command        // Top Command
	Help     IHelpBuilder      // build help
	groups   map[string]string // group name to desc
	root     *Command          // Root Command
}

// New create new App
func New() *App {
	return &App{
		Name: filepath.Base(os.Args[0]),
		Help: &HelpBuilder{},
	}
}

func (app *App) AddGroup(name string, desc string) {
	if app.groups == nil {
		app.groups = make(map[string]string)
	}

	app.groups[name] = desc
}

func (app *App) GetGroup(name string) string {
	if app.groups != nil {
		desc := app.groups[name]
		if desc != "" {
			return desc
		}
	}

	return name
}

func (app *App) MaxCommandNameLen() int {
	length := 0
	for _, cmd := range app.Commands {
		if len(cmd.Name) > length {
			length = len(cmd.Name)
		}
	}

	return length
}

// AddCommands add one command
func (app *App) AddCommands(cmds ...*Command) {
	app.Commands = append(app.Commands, cmds...)
}

// AddFlags add one global flag
func (app *App) AddFlags(flags ...*Flag) {
	app.Flags = append(app.Flags, flags...)
}

// Run process cli
func (app *App) Run() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%+v\n", r)
		}
	}()

	app.setup()
	app.build()
}

func (app *App) setup() {

	hasHelp := false
	for _, cmd := range app.Commands {
		if cmd.isCommand("help") {
			hasHelp = true
			break
		}
	}

	if !hasHelp {
		app.Commands = append(app.Commands, app.DefaultHelp())
	}

	app.root = &Command{}
	app.root.Flags = app.Flags
	app.root.Subs = app.Commands

}

func (app *App) build() {
	// -v --verbose
	// -I/usr/include -I=/usr/include -I /usr/include
	// -aux
	// /c
	// /v:{on/off}

	// find command tree
	cmds := make([]*Command, 0)
	args := make([]string, 0, len(os.Args))
	last := app.root
	cmds = append(cmds, last)
	for idx := 1; idx < len(os.Args); idx++ {
		str := os.Args[idx]
		if str[0] == '/' || str[0] == '-' {
			args = append(args, str)
			continue
		}

		sub := last.findSub(str)
		if sub == nil {
			args = append(args, os.Args[idx:]...)
			break
		}

		cmds = append(cmds, sub)
		last = sub
	}

	// merge flags from all commands
	flags := make(map[string]*Flag)
	for _, cmd := range cmds {
		for _, f := range cmd.Flags {
			alias := strings.Split(f.Name, ",")
			for _, n := range alias {
				flags[strings.TrimSpace(n)] = f
			}
		}
	}

	// build options and params
	isHelp := false
	options := make(map[string]*Flag)
	params := make([]string, 0, len(args))

	for idx := 0; idx < len(args); idx++ {
		str := args[idx]
		style, key, value := parseOption(str)
		if style == styleUnknown {
			params = append(params, str)
			continue
		}

		if key == "" {
			// bad key
			continue
		}

		flag := flags[key]
		if flag == nil {
			if key == "help" || key == "h" {
				isHelp = true
				continue
			}
			if style == styleSingle && app.isMultipleShortOptions(flags, key) {
				// add all single
				for _, s := range key {
					options[string(s)] = flags[string(s)]
				}
			} else {
				// like option but not find, process as param
				params = append(params, str)
			}

			continue
		}

		// parse -I /usr/include
		nextIdx := idx + 1
		if value == "" && nextIdx < len(args) {
			nextStr := args[nextIdx]
			if nextStr[0] != '-' && nextStr[0] != '/' {
				value = nextStr
				idx = nextIdx
			}
		}

		if err := flag.addOption(value); err != nil {
			panic(err)
		}

		options[flag.Key()] = flag
	}

	// check flags required
	for _, flag := range flags {
		if err := flag.validate(); err != nil {
			panic(err)
		}
	}

	// remove root
	cmds = cmds[1:]

	firstCmd := ""
	if len(cmds) > 1 {
		firstCmd = cmds[0].Name
	}

	if !isHelp && len(cmds) == 0 && app.root.findSub(AppCommandName) == nil {
		// no command and no action use help
		isHelp = true
	}

	// process help
	if isHelp && firstCmd != "help" {
		//log.Printf("add help")
		helpCmd := app.root.findSub("help")
		cmds = append([]*Command{helpCmd}, cmds...)
	}

	// invoke
	//log.Printf("%+v,%+v,%+v\n", params, cmds, options)
	ctx := newContext(app, params, cmds, options)
	ctx.Next()
}

func parseOption(str string) (int, string, string) {
	ch := str[0]
	if ch != '-' && ch != '/' {
		return styleUnknown, "", ""
	}

	if ch == '/' {
		// window
		return styleWindow, str[1:], ""
	}

	// unix style '-' or '--'
	if len(str) == 1 {
		return styleSingle, "", ""
	}

	style := styleSingle
	option := str[1:]
	if str[1] == '-' {
		style = styleDouble
		option = str[2:]
	}

	// parse -I/usr/include -I=/usr/include /d:{on/off}
	index := strings.IndexAny(option, "=/:")
	if index != -1 {
		return style, option[:index], option[index+1:]
	} else {
		return style, option, ""
	}
}

func (app *App) isMultipleShortOptions(flags map[string]*Flag, key string) bool {
	for _, s := range key {
		if flags[string(s)] == nil {
			return false
		}
	}

	return true
}

// DefaultHelp 默认的help处理
func (app *App) DefaultHelp() *Command {
	return &Command{
		Name: "help",
		Run: func(ctx *Context) {
			builder := app.Help
			if builder == nil {
				builder = &HelpBuilder{}
			}
			builder.Build(ctx)
		},
	}
}
