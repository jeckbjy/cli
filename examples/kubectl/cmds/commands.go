package cmds

import "github.com/jeckbjy/cli"

func GetCommands() []*cli.Command {
	return []*cli.Command{
		cli.NewCmd("create", "Beginner", onCmdCreate),
		cli.NewCmd("expose", "Beginner", onCmdExpose),
		cli.NewCmd("run", "Beginner", onCmdRun),
		cli.NewCmd("set", "Beginner", onCmdSet),
		cli.NewCmd("explain", "Intermediate", onCmdExplain),
		cli.NewCmd("get", "Intermediate", onCmdGet),
		cli.NewCmd("rollout", "Deploy", onCmdRollout),
		cli.NewCmd("certificate", "Cluster", onCmdCertificate),
		cli.NewCmd("describe", "Debugging", onCmdDescribe),
		cli.NewCmd("diff", "Advanced", onCmdDiff),
		cli.NewCmd("label", "Settings", onCmdLabel),
		cli.NewCmd("api-resources", "Other", onCmdAPIResources),
	}
}
