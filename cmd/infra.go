package cmd

import (
	"fmt"
	"github.com/controlplaneio/simulator-standalone/pkg/simulator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `create`,
		Short: "Runs terraform to create the required infrastructure for scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			tfDir := viper.GetString("tf-dir")
			return simulator.Create(tfDir)
		},
	}

	return cmd
}

func newStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `status`,
		Short: "Gets the status of the infrastructure",
		RunE: func(cmd *cobra.Command, args []string) error {
			tfDir := viper.GetString("tf-dir")
			tfo, err := simulator.Status(tfDir)

			fmt.Println(tfo)

			return err
		},
	}

	return cmd
}

func newDestroyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `destroy`,
		Short: "Tears down the infrastructure created for scenarios",
		RunE: func(cmd *cobra.Command, args []string) error {
			tfDir := viper.GetString("tf-dir")
			return simulator.Destroy(tfDir)
		},
	}

	return cmd
}

func newInfraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           `infra <subcommand>`,
		Short:         "Interact with AWS to create, query and destroy the required infrastructure for scenarios",
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	cmd.AddCommand(newCreateCommand())
	cmd.AddCommand(newStatusCommand())
	cmd.AddCommand(newDestroyCommand())

	return cmd
}
