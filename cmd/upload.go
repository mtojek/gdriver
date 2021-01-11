package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mtojek/gdriver/internal/auth"
	"github.com/mtojek/gdriver/internal/driveext"
	"github.com/mtojek/gdriver/internal/upload"
)

func setupUploadCommand() *cobra.Command {
	uploadCmd := &cobra.Command{
		Use:          "upload [folderID]",
		Short:        "Upload files",
		Long:         "Use upload subcommand to upload files to Google Drive.",
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

			sourceDir, _ := cmd.Flags().GetString("source")
			selectionMode, _ := cmd.Flags().GetBool("select")

			driveService, err := driveext.NewService()
			if err != nil {
				return errors.Wrap(err, "initializing drive service failed")
			}

			err = upload.Files(driveService, upload.FilesOptions{
				FolderID:      folderID,
				SourceDir:     sourceDir,
				SelectionMode: selectionMode,
			})
			if err != nil {
				return errors.Wrap(err, "uploading files failed")
			}

			fmt.Println("Done")
			return nil
		},
	}
	uploadCmd.Flags().String("source", ".", "Source folder with resources to upload")
	uploadCmd.Flags().Bool("select", false, "Select files to upload")
	return uploadCmd
}
