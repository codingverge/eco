package main

import (
	"os"

	"github.com/codingverge/axon/config"
	"github.com/codingverge/axon/driver"
	"github.com/spf13/cobra"
)

func main() {
	Execute()
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server",
	Long:  `server -c example/server/config.yaml`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		drv, err := driver.New(ctx, cmd.OutOrStderr(), nil, []config.Option{config.WithFlags(cmd.Flags())})

		if err != nil {
			return err
		}
		return drv.RunE(ctx)
	},
}

func init() {
	// registering flags to config
	config.RegisterConfigFlag(serverCmd.PersistentFlags(), []string{})
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := serverCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
