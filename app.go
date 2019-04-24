package cli

import (
	"fmt"
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

// colors
// https://github.com/logrusorgru/aurora
// https://github.com/mgutz/ansi
// https://medium.com/@inhereat/terminal-color-rendering-tool-library-support-8-16-colors-256-colors-by-golang-a68fb8deee86
// https://github.com/gookit/color

// App build a git style cli
type App struct {
	Name      string
	help      IHelp             // custom help
	root      *Command          // Root Command
	groups    map[string]string // group name to desc
	languages map[string]string // language map
}

// New create new App
func New() *App {
	return &App{
		Name: filepath.Base(os.Args[0]),
		root: &Command{},
	}
}

// NewWithHelp create new app with help
func NewWithHelp(help IHelp) *App {
	return &App{
		Name: filepath.Base(os.Args[0]),
		help: help,
		root: &Command{},
	}
}

func (app *App) Root() *Command {
	return app.root
}

func (app *App) SetLanguage(langs map[string]string) {
	app.languages = langs
}

func (app *App) SetGroups(groups map[string]string) {
	app.groups = groups
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

// AddCommands add one command
func (app *App) AddCommands(cmds []*Command) {
	app.root.Subs = append(app.root.Subs, cmds...)
}

// AddFlags add one global flag
func (app *App) AddFlags(flags []*Flag) {
	app.root.Flags = append(app.root.Flags, flags...)
}

// AddHeader add header to root command
func (app *App) AddHeader(data string) {
	app.root.Header += data
}

// AddFooter add footer to root command
func (app *App) AddFooter(data string) {
	app.root.Footer += data
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
	app.Name = strings.TrimSpace(app.Name)
	if app.help == nil {
		app.help = &Help{}
	}

	if app.root.findSub("help") == nil {
		app.root.AddSub(&Command{
			Name: "help",
			Run: func(ctx *Context) {
				ctx.App().help.Build(ctx)
			},
		})
	}
}

func (app *App) build() {
	// -v --verbose
	// -I/usr/include -I=/usr/include -I /usr/include
	// -aux
	// /c
	// /v:{on/off}

	// find command list and args
	cmds, args := app.buildCommands()
	// merge flags from all commands
	flags := app.buildAllFlags(cmds)

	// build options and params
	isHelp := false
	options := make(map[string]*Flag)
	params := make([]string, 0, len(args))

	for idx := 0; idx < len(args); idx++ {
		str := args[idx]
		style, key, value := app.parseOption(str)
		if style == styleUnknown {
			params = append(params, str)
			continue
		}

		if key == "" {
			// bad key
			continue
		}

		if key == "help" || key == "h" {
			isHelp = true
			continue
		}

		var flag *Flag

		// check multiple short option
		if style == styleSingle && len(key) > 1 {
			for i := 0; i < len(key); i++ {
				ch := key[i]
				st := string(ch)
				flag = flags[st]
				if flag == nil {
					panic(fmt.Errorf("Unknown option %+v", st))
				}
				options[flag.Name] = flag
			}
		} else {
			flag = flags[key]
			if flag == nil {
				panic(fmt.Errorf("Unknown option %+v", key))
			}
			options[flag.Name] = flag
		}

		// need check flag has params?

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
	}

	// check flags required
	for _, flag := range flags {
		if err := flag.validate(); err != nil {
			panic(err)
		}
	}

	// remove root
	cmds = cmds[1:]

	if !isHelp && len(cmds) == 0 && app.root.findSub(AppCommandName) == nil {
		// no command and no action use help
		isHelp = true
	}

	if len(cmds) > 0 && cmds[0].Name == "help" {
		isHelp = true
		cmds = cmds[1:]
	}

	ctx := newContext(app, params, cmds, options)

	if isHelp {
		app.root.findSub("help").Run(ctx)
	} else {
		ctx.Next()
	}
}

func (app *App) buildCommands() ([]*Command, []string) {
	cmds := []*Command{app.root}
	args := make([]string, 0, len(os.Args))
	last := app.root
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
		if sub.Name != "help" {
			last = sub
		}
	}

	return cmds, args
}

func (app *App) buildAllFlags(cmds []*Command) map[string]*Flag {
	flags := make(map[string]*Flag)
	for _, cmd := range cmds {
		for _, f := range cmd.Flags {
			if f.Name != "" {
				flags[f.Name] = f
			}

			if f.Short != "" {
				flags[f.Name] = f
			}
		}
	}

	return flags
}

func (app *App) parseOption(str string) (int, string, string) {
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
	}

	return style, option, ""
}

// Translate translate to locale
func (app *App) Translate(name string, alternative string) string {
	if name != "" && name[0] != '$' {
		return name
	}

	if app.languages == nil {
		return ""
	}

	var key string
	if name != "" {
		key = name[1:]
	} else {
		key = alternative
	}

	if value, ok := app.languages[key]; ok {
		return value
	}

	return ""
}
