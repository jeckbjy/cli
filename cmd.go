package cli

// Action command callback
type Action func(ctx *Context)

// CommandGroup command set
type CommandGroup struct {
	Name string
	Cmds []*Command
}

// Command cli command tree
type Command struct {
	Name   string
	Group  string
	Short  string
	Long   string
	Header string
	Footer string
	Run    Action
	Flags  []*Flag
	Subs   []*Command
	Alias  []string
}

// NewCmd create command
func NewCmd(name string, group string, action Action) *Command {
	return &Command{
		Name:  name,
		Group: group,
		Run:   action,
	}
}

func (cmd *Command) isCommand(name string) bool {
	if cmd.Name == name {
		return true
	}

	for _, value := range cmd.Alias {
		if value == name {
			return true
		}
	}

	return false
}

func (cmd *Command) findSub(name string) *Command {
	for _, sub := range cmd.Subs {
		if sub.isCommand(name) {
			return sub
		}
	}

	return nil
}

func (cmd *Command) AddSub(sub *Command) {
	cmd.Subs = append(cmd.Subs, sub)
}

// MaxSubNameLen return sub command name length
func (cmd *Command) MaxSubNameLen() int {
	length := 0
	for _, sub := range cmd.Subs {
		if len(sub.Name) > length {
			length = len(sub.Name)
		}
	}

	return length
}

func (cmd *Command) SubCommandGroups() []*CommandGroup {
	groupList := make([]*CommandGroup, 0)
	groupMap := make(map[string]*CommandGroup)
	global := &CommandGroup{}
	for _, sub := range cmd.Subs {
		var group *CommandGroup
		// ignore help
		if sub.Name == "help" {
			continue
		}

		if sub.Group == "" {
			group = global
		} else if group = groupMap[sub.Group]; group == nil {
			group = &CommandGroup{Name: sub.Group}
			groupMap[sub.Group] = group
			groupList = append(groupList, group)
		}

		group.Cmds = append(group.Cmds, sub)
	}

	if len(global.Cmds) != 0 {
		groupList = append(groupList, global)
	}

	return groupList
}
