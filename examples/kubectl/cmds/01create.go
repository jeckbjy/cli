package cmds

import "github.com/jeckbjy/cli"

var createLongDesc = `
Create a resource from a file or from stdin.

JSON and YAML formats are accepted.

Examples:
  # Create a pod using the data in pod.json.
  kubectl create -f ./pod.json

  # Create a pod based on the JSON passed into stdin.
  cat pod.json | kubectl create -f -

  # Edit the data in docker-registry.yaml in JSON then create the resource using the edited data.
  kubectl create -f docker-registry.yaml --edit -o json

Available Commands:
  clusterrole         Create a ClusterRole.
  clusterrolebinding  为一个指定的 ClusterRole 创建一个 ClusterRoleBinding
  configmap           从本地 file, directory 或者 literal value 创建一个 configmap
  deployment          创建一个指定名称的 deployment.
  job                 Create a job with the specified name.
  namespace           创建一个指定名称的 namespace
  poddisruptionbudget 创建一个指定名称的 pod disruption budget.
  priorityclass       Create a priorityclass with the specified name.
  quota               创建一个指定名称的 quota.
  role                Create a role with single rule.
  rolebinding         为一个指定的 Role 或者 ClusterRole创建一个 RoleBinding
  secret              使用指定的 subcommand 创建一个 secret
  service             使用指定的 subcommand 创建一个 service.
  serviceaccount      创建一个指定名称的 service account

Options:
      --allow-missing-template-keys=true: If true, ignore any errors in templates when a field or map key is missing in
the template. Only applies to golang and jsonpath output formats.
      --dry-run=false: If true, only print the object that would be sent, without sending it.
      --edit=false: Edit the API resource before creating
  -f, --filename=[]: Filename, directory, or URL to files to use to create the resource
  -o, --output='': Output format. One of:
json|yaml|name|go-template|go-template-file|template|templatefile|jsonpath|jsonpath-file.
      --raw='': Raw URI to POST to the server.  Uses the transport specified by the kubeconfig file.
      --record=false: Record current kubectl command in the resource annotation. If set to false, do not record the
command. If set to true, record the command. If not set, default to updating the existing annotation value only if one
already exists.
  -R, --recursive=false: Process the directory used in -f, --filename recursively. Useful when you want to manage
related manifests organized within the same directory.
      --save-config=false: If true, the configuration of current object will be saved in its annotation. Otherwise, the
annotation will be unchanged. This flag is useful when you want to perform kubectl apply on this object in the future.
  -l, --selector='': Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)
      --template='': Template string or path to template file to use when -o=go-template, -o=go-template-file. The
template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview].
      --validate=true: If true, use a schema to validate the input before sending it
      --windows-line-endings=false: Only relevant if --edit=true. Defaults to the line ending native to your platform.

Usage:
  kubectl create -f FILENAME [options]

Use "kubectl <command> --help" for more information about a given command.
Use "kubectl options" for a list of global command-line options (applies to all commands).
`

func onCmdCreate(ctx *cli.Context) {
	var flags struct {
		AllowMissingTemplateKeys bool   `cli:"allow-missing-template-keys"`
		DryRun                   bool   `cli:"dry-run"`
		Edit                     bool   `cli:"edit"`
		Filename                 string `cli:"filename"`
		Output                   string `cli:"output"`
		Raw                      string `cli:"raw"`
		Record                   bool
		Recursive                bool
		SaveConfig               bool
		Selector                 string
	}

	ctx.Bind(&flags)

	if ctx.NArg() < 0 {
		ctx.Abort()
	}

	switch ctx.Arg(0) {
	case "clusterrole":
	case "clusterrolebinding":
	default:
		ctx.Abort()
	}
}
