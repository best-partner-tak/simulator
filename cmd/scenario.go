package cmd

import (
	"fmt"
	"github.com/controlplaneio/simulator-standalone/pkg/scenario"
	"github.com/controlplaneio/simulator-standalone/pkg/simulator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newScenarioListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `list`,
		Short: "Lists available scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestPath := viper.GetString("scenarios-dir")
			manifest, err := scenario.LoadManifest(manifestPath)

			if err != nil {
				return err
			}

			fmt.Println("Available scenarios:")
			for _, s := range manifest.Scenarios {
				fmt.Println("ID: " + s.Id + ", Name: " + s.DisplayName)
			}

			return nil
		},
	}

	return cmd
}

func newScenarioLaunchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `launch <id>`,
		Short: "Launches a scenario",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("scenario id is required")
			}

			bucket := viper.GetString("bucket")
			tfDir := viper.GetString("tf-dir")
			scenariosDir := viper.GetString("scenarios-dir")
			scenarioID := args[0]

			return simulator.Launch(tfDir, scenariosDir, bucket, scenarioID)
		},
	}

	return cmd
}

func newScenarioCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           `scenario <subcommand>`,
		Short:         "Interact with scenarios",
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	cmd.AddCommand(newScenarioListCommand())
	cmd.AddCommand(newScenarioLaunchCommand())

	return cmd
}
