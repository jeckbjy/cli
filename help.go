package cli

import (
	"fmt"
	"strings"
)

// IHelp build help
type IHelp interface {
	Build(ctx *Context)
}

// Help build help
// Default template:
// {{header}}
// Usage: {{app}} <command> <args> [options]
// Commands:
//   {{commands}}
// Options:
//   {{options}}
// {{footer}}
type Help struct {
	indent string
	last   string
	size   int
}

// SetIndent set indent
func (h *Help) SetIndent(indent string) {
	h.indent = indent
}

// Write print string without indent, and auto add \n and ignore ""
func (h *Help) Write(format string, args ...interface{}) {
	content := fmt.Sprintf(format, args...)
	if content == "" {
		return
	}

	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	fmt.Printf(content)

	h.size += len(content)
	if len(content) > 2 {
		h.last = content
	} else {
		h.last += content
	}
}

// WriteIndent print content with prefix, auto split multi lines and add space
// -v, --verbose               Noisy logging, including all shell commands executed.
//                             If used with --help, shows hidden options.
func (h *Help) WriteIndent(prefix string, format string, args ...interface{}) {
	if prefix == "" {
		h.Write(format, args...)
		return
	}

	content := fmt.Sprintf(format, args...)
	if content == "" {
		h.Write(prefix)
		return
	}

	// TODO: display According to the size of the console
	space := len(prefix)
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		str := fmt.Sprintf("%s%s\n", prefix, line)
		h.Write(str)

		if i == 0 {
			prefix = fmt.Sprintf("%*s", space, "")
		}
	}
}

// WriteDivide 输出分隔行
func (h *Help) WriteDivide() {
	if h.size == 0 {
		return
	}

	if strings.HasSuffix(h.last, "\n\n") {
		return
	}

	h.Write("\n")
}

// Indent return indent count
func (h *Help) Indent(data string, len int) string {
	return fmt.Sprintf("%s%-*s", h.indent, len, data)
}

// Build default build help
func (h *Help) Build(ctx *Context) {
	if h.indent == "" {
		h.indent = "  "
	}

	target := h.GetTargetCommand(ctx)

	// write header
	h.Write(target.Header)
	h.WriteDivide()
	h.WriteSubCommands(ctx, target)
	h.WriteDivide()
	h.WriteOptions(ctx, target, "Available options:", true)
	h.WriteDivide()
	h.WriteUsage(ctx)
	h.WriteDivide()
	h.Write(target.Footer)
}

func (h *Help) GetTargetCommand(ctx *Context) *Command {
	cmds := ctx.CommandList()
	if len(cmds) == 0 {
		return ctx.App().Root()
	}
	return cmds[len(cmds)-1]
}

func (h *Help) WriteUsage(ctx *Context) {
	app := ctx.App()
	cmds := ctx.CommandList()

	if len(cmds) == 0 {
		h.Write("Usage: %s [<args>] [<options>]", app.Name)
	} else {
		builder := strings.Builder{}

		for _, c := range cmds {
			if builder.Len() > 0 {
				builder.WriteString(" ")
			}

			builder.WriteString(c.Name)
		}

		cmdsName := builder.String()

		h.Write("Usage: %s %s [<args>] [<options>]", app.Name, cmdsName)
	}
}

func (h *Help) WriteSubCommands(ctx *Context, cmd *Command) {
	if len(cmd.Subs) == 0 {
		return
	}

	app := ctx.App()

	groups := cmd.SubCommandGroups()
	if len(groups) == 0 {
		return
	}

	maxNameLen := cmd.MaxSubNameLen() + len(h.indent)

	for _, g := range groups {
		h.WriteDivide()

		if g.Name != "" {
			h.Write(app.GetGroup(g.Name))
		} else {
			h.Write("Available commands:")
		}

		for _, c := range g.Cmds {
			h.WriteIndent(h.Indent(c.Name, maxNameLen), app.Translate(c.Short, c.Name+"_s"))
		}
	}
}

func (h *Help) WriteOptions(ctx *Context, cmd *Command, head string, hasIndent bool) {
	if len(cmd.Flags) == 0 {
		return
	}

	maxLen := 0
	for _, f := range cmd.Flags {
		fname := f.FullName()
		if fname != "" && len(fname) > maxLen {
			maxLen = len(fname)
		}
	}

	indent := ""
	if hasIndent {
		indent = h.indent
	}

	if head != "" {
		h.Write(head)
	}

	for _, f := range cmd.Flags {
		prefix := ""
		if f.Short != "" {
			prefix = fmt.Sprintf("%s-%s,--%-*s", indent, f.Short, maxLen, f.FullName())
		} else {
			prefix = fmt.Sprintf("%s   --%-*s", indent, maxLen, f.FullName())
		}

		h.WriteIndent(prefix, f.Usage)
	}
}
