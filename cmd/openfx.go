package cmd

import (
	"os"

	"github.com/keti-openfx/openfx-cli/cmd/function"
	"github.com/keti-openfx/openfx-cli/cmd/runtime"
	"github.com/spf13/cobra"
)

var openfxCmd = &cobra.Command{
	Use:   "openfx-cli",
	Short: "Manage Openfx",
	Long: `
	Manage Openfx functions from the command line interface
	`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {

	openfxCmd.SetUsageTemplate(usageTemplate)
	openfxCmd.SetHelpTemplate(helpTemplate)

	openfxCmd.AddCommand(versionCmd)
	openfxCmd.AddCommand(function.FunctionCmd)
	openfxCmd.AddCommand(runtime.RuntimeCmd)
}

func Execute() {
	if err := openfxCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if .IsAvailableCommand }}
  {{rpad .NameAndAliases 20}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}

`

var helpTemplate = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}
{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
