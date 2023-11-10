package cli

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nikolaymatrosov/sls-rosetta/internal/examples"
	"github.com/nikolaymatrosov/sls-rosetta/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const cfg = "config"

func RootCmd(cmd *cobra.Command, args []string) error {
	config, err := examples.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	m := ui.NewViewModel(config)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return nil
}
func RootPreRun(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}
	viper.AutomaticEnv()
	viper.SetEnvPrefix("rosetta")

	config := examples.NewConfig()
	cfgPath := viper.GetString(cfg)

	if _, err := os.Stat(cfgPath); errors.Is(err, os.ErrNotExist) {
		if err := config.Fetch(examples.DefaultConfigUrl()); err != nil {
			return err
		}
		cmd.SetContext(examples.ContextWith(cmd.Context(), config))
		return nil
	}

	cfgFile, err := os.ReadFile(cfgPath)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(cfgFile, config); err != nil {
		return err
	}
	cmd.SetContext(examples.ContextWith(cmd.Context(), config))
	return nil
}
