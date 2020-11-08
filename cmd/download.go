package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mtojek/gdriver/internal/auth"
	"github.com/mtojek/gdriver/internal/download"
	"github.com/mtojek/gdriver/internal/driveext"
)

func setupDownloadCommand() *cobra.Command {
	downloadCmd := &cobra.Command{
		Use:          "download [folderID]",
		Short:        "Download files",
		Long:         "Use download subcommand to download files from Google Drive.",
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

			outputDir, _ := cmd.Flags().GetString("output")
			selectionMode, _ := cmd.Flags().GetBool("select")

			driveService, err := driveext.NewService()
			if err != nil {
				return errors.Wrap(err, "initializing drive service failed")
			}

			err = download.Files(driveService, download.FilesOptions{
				FolderID:      folderID,
				OutputDir:     outputDir,
				SelectionMode: selectionMode,
			})
			if err != nil {
				return errors.Wrap(err, "downloading files failed")
			}

			fmt.Println("Done")
			return nil
		},
	}
	downloadCmd.Flags().String("output", ".", "Output folder for downloaded resources")
	downloadCmd.Flags().Bool("select", false, "Select files to download")
	return downloadCmd
}
