package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/stjudecloud/msgenctl/internal"
)

var rootCmd = &cobra.Command{
	Version:      internal.Version,
	Use:          "msgenctl",
	Short:        "Query and send commands to Microsoft Genomics",
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.InitDefaultVersionFlag()

	persistentFlags := rootCmd.PersistentFlags()

	persistentFlags.String("base-url", "", "Microsoft Genomics API base URL")
	rootCmd.MarkPersistentFlagRequired("base-url")

	persistentFlags.String("access-key", "", "Microsoft Genomics API access key")
	rootCmd.MarkPersistentFlagRequired("access-key")
}
