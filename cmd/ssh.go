package cmd

import (
	"fmt"
	"github.com/controlplaneio/simulator-standalone/pkg/simulator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newSSHConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `config`,
		Short: "Prints the stanzas to add to ssh config to connect to your cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			scenariosDir := viper.GetString("scenarios-dir")
			bucket := viper.GetString("bucket")
			tfDir := viper.GetString("tf-dir")
			cfg, err := simulator.Config(tfDir, scenariosDir, bucket)
			if err != nil {
				return err
			}

			fmt.Println(*cfg)

			return nil
		},
	}

	return cmd
}

func newSSHAttackCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `attack`,
		Short: "Connect to an attack container to complete the scenario",
		RunE: func(cmd *cobra.Command, args []string) error {
			bucket := viper.GetString("bucket")
			tfDir := viper.GetString("tf-dir")

			return simulator.Attack(tfDir, bucket)
		},
	}

	return cmd

}

func newSSHCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           `ssh <command>`,
		Short:         "Interact with the cluster",
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	cmd.AddCommand(newSSHConfigCommand())
	cmd.AddCommand(newSSHAttackCommand())

	return cmd
}
