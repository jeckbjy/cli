package main

import (
	"github.com/jeckbjy/cli"
	"github.com/jeckbjy/cli/examples/kubectl/cmds"
)

func main() {
	app := cli.New()
	app.SetLanguage(cmds.Langs)

	app.AddHeader("kubectl controls the Kubernetes cluster manager.\n")
	app.AddHeader("\n")
	app.AddHeader("Find more information at: https://kubernetes.io/docs/reference/kubectl/overview/\n")
	app.AddHeader("\n")

	app.AddFooter(`Use "kubectl <command> --help" for more information about a given command.`)
	app.AddFooter("\n")
	app.AddFooter(`Use "kubectl options" for a list of global command-line options (applies to all commands).`)
	app.AddFooter("\n")

	app.AddGroup("Beginner", "Basic Commands (Beginner):")
	app.AddGroup("Intermediate", "Basic Commands (Intermediate):")
	app.AddGroup("Deploy", "Deploy Commands:")
	app.AddGroup("Cluster", "Cluster Management Commands:")
	app.AddGroup("Debugging", "Troubleshooting and Debugging Commands:")
	app.AddGroup("Advanced", "Advanced Commands:")
	app.AddGroup("Settings", "Settings Commands:")
	app.AddGroup("Other", "Other Commands:")

	app.AddCommands(cmds.GetCommands())
	app.Run()
}
