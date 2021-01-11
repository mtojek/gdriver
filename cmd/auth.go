package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mtojek/gdriver/internal/auth"
)

func setupAuthCommand() *cobra.Command {
	authCmd := &cobra.Command{
		Use:          "auth",
		Short:        "Authenticate Google account",
		Long:         "Use auth subcommand to authenticate with Google Drive API.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			newCredentialsFile, _ := cmd.Flags().GetString("import-credentials")
			readOnlyScope, _ := cmd.Flags().GetBool("read-only")
			err := auth.Authenticate(newCredentialsFile, readOnlyScope)
			if err != nil {
				return errors.Wrap(err, "authentication failed")
			}
			return nil
		},
	}
	authCmd.Flags().String("import-credentials", "", "Client credentials file (for Google Drive API)")
	authCmd.Flags().Bool("read-only", true, "Read-only Drive scope")
	return authCmd
}
