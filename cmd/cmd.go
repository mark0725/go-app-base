package cmd

import (
	"fmt"

	// "github.com/containerd/console"
	"github.com/spf13/cobra"
)

func NewCLI(options *CmdOptions) *cobra.Command {
	cobra.EnableCommandSorting = false

	// if runtime.GOOS == "windows" {
	// 	console.ConsoleFromFile(os.Stdin) //nolint:errcheck
	// }

	rootCmd := &cobra.Command{
		Use:           options.AppName,
		Short:         options.Description,
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Run: func(cmd *cobra.Command, args []string) {
			if version, _ := cmd.Flags().GetBool("version"); version {
				fmt.Printf("%s version is %s\n", options.AppName, options.AppVersion)
				return
			}

			cmd.Print(cmd.UsageString())
		},
	}

	rootCmd.Flags().BoolP("version", "v", false, "Show version information")

	return rootCmd
}
