package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mtojek/gdriver/internal/auth"
	"github.com/mtojek/gdriver/internal/download"
)

func Root() *cobra.Command {
	authCmd := &cobra.Command{
		Use:          "auth",
		Short:        "Authenticate Google account",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			newCredentialsFile, _ := cmd.Flags().GetString("import-credentials")
			err := auth.Authenticate(newCredentialsFile)
			if err != nil {
				return errors.Wrap(err, "authentication failed")
			}
			return nil
		},
	}
	authCmd.Flags().String("import-credentials", "", "Client credentials file (for Google Drive API)")

	downloadCmd := &cobra.Command{
		Use:          "download [folderID]",
		Short:        "Download files",
		SilenceUsage: true,
		Args:         cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := auth.Verify()
			if err != nil {
				return errors.Wrap(err, "auth verification failed")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var folderID string
			if len(args) > 0 {
				folderID = args[0]
			}

			outputPath, _ := cmd.Flags().GetString("output")
			err := download.Files(folderID, outputPath)
			if err != nil {
				return errors.Wrap(err, "downloading files failed")
			}
			return nil
		},
	}
	downloadCmd.Flags().String("output", ".", "Output folder for downloaded resources")

	rootCmd := &cobra.Command{
		Use:   "gdriver",
		Short: "Download large files from Google Drive",
		Long:  "Use gdriver to download large files from Google Drive",
	}
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(downloadCmd)
	return rootCmd
}
