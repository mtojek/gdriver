package cmd

import (
	"github.com/spf13/cobra"
)

func Root() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gdriver",
		Short: "Download and upload large files to Google Drive",
		Long:  "Use gdriver to download and upload large files to Google Drive.",
	}
	rootCmd.AddCommand(setupAuthCommand())
	rootCmd.AddCommand(setupCheckCommand())
	rootCmd.AddCommand(setupDownloadCommand())
	rootCmd.AddCommand(setupUploadCommand())
	return rootCmd
}
