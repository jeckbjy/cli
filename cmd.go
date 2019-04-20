package cli

// Action command callback
type Action func(ctx *Context)

// Command cli command tree
type Command struct {
	Name  string
	Group string
	Short string
	Long  string
	Run   Action
	Flags []*Flag
	Subs  []*Command
	Alias []string
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
