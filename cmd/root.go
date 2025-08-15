package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pull2push",
	Short: "A live streaming proxy service",
	Long: `pull2push Live Client is a service that handles live streaming proxying and management.
It provides APIs for creating, managing, and monitoring live streams.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},

	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		// 清理资源
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file path")
	rootCmd.MarkPersistentFlagRequired("config")
	rootCmd.AddCommand(serveCmd)
}
