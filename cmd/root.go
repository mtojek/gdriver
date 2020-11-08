package cmd

import (
	"github.com/spf13/cobra"
)

func Root() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gdriver",
		Short: "Download large files from Google Drive",
		Long:  "Use gdriver to download large files from Google Drive.",
	}
	rootCmd.AddCommand(setupAuthCommand())
	rootCmd.AddCommand(setupCheckCommand())
	rootCmd.AddCommand(setupDownloadCommand())
	return rootCmd
}
