package cli

import (
	"fmt"
	"strings"
)

type CommandGroup struct {
	Name     string
	Desc     string
	Commands []*Command
}

type IHelpBuilder interface {
	Build(ctx *Context)
}

type HelpBuilder struct {
	indent string
	prefix string
	depth  int
	locale map[string]string // default local language
	header strings.Builder
	footer strings.Builder
}

func (b *HelpBuilder) Header() string {
	return b.header.String()
}

func (b *HelpBuilder) Footer() string {
	return b.footer.String()
}

// AddHeader append desc to header and auto add \n
func (b *HelpBuilder) AddHeader(desc string) {
	b.header.WriteString(desc)
	if !strings.HasSuffix(desc, "\n") {
		b.header.WriteString("\n")
	}
}

// AddFooter append desc to footer and auto add \n
func (b *HelpBuilder) AddFooter(desc string) {
	b.footer.WriteString(desc)
	if !strings.HasSuffix(desc, "\n") {
		b.footer.WriteString("\n")
	}
}

func (b *HelpBuilder) SetLocale(l map[string]string) {
	b.locale = l
}

func (b *HelpBuilder) Push() {
	b.depth++
	b.prefix += b.indent
}

func (b *HelpBuilder) Pop() {
	if b.depth > 0 {
		b.depth--
		b.prefix = b.prefix[:len(b.prefix)-len(b.indent)]
	}
}

// Write auto split \n and auto add \n and auto add prefix
func (b *HelpBuilder) Write(format string, args ...interface{}) {
	data := fmt.Sprintf(format, args...)

	if b.depth > 0 {
		lines := strings.Split(data, "\n")
		for _, line := range lines {
			fmt.Printf("%s%s\n", b.prefix, line)
		}
	} else {
		fmt.Printf("%s", data)
		if !strings.HasSuffix(data, "\n") {
			fmt.Println()
		}
	}
}

func (b *HelpBuilder) Build(ctx *Context) {
	if b.indent == "" {
		b.indent = "  "
	}

	// build
	//fmt.Printf("build help\n")
	commands := ctx.CommandList()
	if len(commands) == 1 {
		// global help
		b.BuildApp(ctx)
	} else {
		// top command
		b.BuildCmd(ctx, commands[1])
	}
}

func (b *HelpBuilder) Locale(text string) string {
	if text == "" {
		return ""
	}

	if text[0] != '$' {
		return text
	}

	// TODO: load locale from file
	if b.locale != nil {
		return b.locale[text[1:]]
	}

	// not found
	return ""
}

func (b *HelpBuilder) LocaleShort(cmd string, text string) string {
	if text == "" {
		return b.Locale(fmt.Sprintf("$%s_s", cmd))
	}

	return b.Locale(text)
}

func (b *HelpBuilder) LocaleLong(cmd string, text string) string {
	if text == "" {
		return b.Locale(fmt.Sprintf("$%s_l", cmd))
	}

	return b.Locale(text)
}

func (b *HelpBuilder) GetCommandGroup(ctx *Context) []*CommandGroup {
	groupList := make([]*CommandGroup, 0)
	groupMap := make(map[string]*CommandGroup)

	app := ctx.app

	for _, cmd := range app.Commands {
		// ignore help
		if cmd.Name == "help" {
			continue
		}

		name := cmd.Group

		group := groupMap[name]
		if group == nil {
			group = &CommandGroup{Name: name, Desc: app.GetGroup(name)}
			groupMap[name] = group
			groupList = append(groupList, group)
		}

		group.Commands = append(group.Commands, cmd)
	}

	return groupList
}

func (b *HelpBuilder) BuildApp(ctx *Context) {
	app := ctx.App()
	b.Write(b.Header())

	// build command group
	maxLen := app.MaxCommandNameLen()
	groups := b.GetCommandGroup(ctx)
	for _, group := range groups {
		if group.Name != "" {
			b.Write(group.Desc)
			b.Push()
			for _, cmd := range group.Commands {
				b.Write("%-*s%s%s", maxLen, cmd.Name, b.indent, b.LocaleShort(cmd.Name, cmd.Short))
			}
			b.Pop()
		} else {
			for _, cmd := range group.Commands {
				b.Write("%-*s%s%s", maxLen, cmd.Name, b.indent, b.LocaleShort(cmd.Name, cmd.Short))
			}
		}
		b.Write("")
	}

	// build usage
	b.Write("Usage:")
	b.Write("%s%s [flags][options]", b.indent, app.Name)
	b.Write("")

	b.Write(b.Footer())
}

func (b *HelpBuilder) BuildCmd(ctx *Context, cmd *Command) {

}
