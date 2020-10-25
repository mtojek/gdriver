package cmd

import (
	"github.com/spf13/cobra"
)

func Root() *cobra.Command {
	authCmd := &cobra.Command{
		Use: "auth",
		Short: "Authenticate Google user account",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	downloadCmd := &cobra.Command{
		Use: "download [folderID]",
		Short: "Download files",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	downloadCmd.Flags().String("output", ".", "Output folder for downloaded resources")

	rootCmd := &cobra.Command{
		Use:   "gdriver",
		Short: "Manage large files in Google Drive",
		Long:  "Use gdriver to manage large files in Google Drive",
	}
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(downloadCmd)
	return rootCmd
}
