package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mtojek/gdriver/internal/auth"
	"github.com/mtojek/gdriver/internal/check"
	"github.com/mtojek/gdriver/internal/driveext"
)

func setupCheckCommand() *cobra.Command {
	checkCmd := &cobra.Command{
		Use:          "check [folderID]",
		Short:        "Check files",
		Long:         "Use check subcommand to verify files downloaded from Google Drive.",
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

			targetDir, _ := cmd.Flags().GetString("target")

			driveService, err := driveext.NewService()
			if err != nil {
				return errors.Wrap(err, "initializing drive service failed")
			}

			err = check.Files(driveService, check.FilesOptions{
				FolderID:  folderID,
				TargetDir: targetDir,
			})
			if err != nil {
				return errors.Wrap(err, "checking files failed")
			}

			fmt.Println("Done")
			return nil
		},
	}
	checkCmd.Flags().String("target", ".", "Target directory with downloaded resources")
	return checkCmd
}
