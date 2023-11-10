package cli

import (
	"context"

	"github.com/spf13/cobra"
)

// Execute is the command line applications entry function
func Execute() error {
	rootCmd := &cobra.Command{
		Version:           "v0.0.1",
		Use:               "rosetta",
		Long:              "Rosetta is a tool for generating serverless functions from templates.",
		Example:           "rosetta",
		RunE:              RootCmd,
		PersistentPreRunE: RootPreRun,
	}

	var cfgFile string

	rootCmd.PersistentFlags().StringVar(&cfgFile, cfg, "", "config file (default is $HOME/.sls-rosetta.yaml)")

	rootCmd.AddCommand(initialize())
	return rootCmd.ExecuteContext(context.Background())
}

func initialize() *cobra.Command {
	init := &cobra.Command{
		Use:     "initialize",
		Short:   "init rosetta configuration file.",
		Long:    "init provisions a rosetta configuration file.",
		Example: "rosetta init",
		Aliases: []string{"i", "init"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return init
}
