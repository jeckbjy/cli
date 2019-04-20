package main

import (
	"github.com/jeckbjy/cli"
	"github.com/jeckbjy/cli/examples/kubectl/cmds"
)

func main() {
	app := cli.New()
	help := &cli.HelpBuilder{}
	help.SetLocale(cmds.Locales)
	help.AddHeader("kubectl controls the Kubernetes cluster manager.")
	help.AddHeader("")
	help.AddHeader("Find more information at: https://kubernetes.io/docs/reference/kubectl/overview/")
	help.AddHeader("")

	help.AddFooter(`Use "kubectl <command> --help" for more information about a given command.`)
	help.AddFooter(`Use "kubectl options" for a list of global command-line options (applies to all commands).`)

	app.AddGroup("Beginner", "Basic Commands (Beginner):")
	app.AddGroup("Intermediate", "Basic Commands (Intermediate):")
	app.AddGroup("Deploy", "Deploy Commands:")
	app.AddGroup("Cluster", "Cluster Management Commands:")
	app.AddGroup("Debugging", "Troubleshooting and Debugging Commands:")
	app.AddGroup("Advanced", "Advanced Commands:")
	app.AddGroup("Settings", "Settings Commands:")
	app.AddGroup("Other", "Other Commands:")

	app.Help = help
	// app.AddFlags(cli.NewFlag())
	app.AddCommands(cmds.GetCommands()...)
	app.Run()
}
