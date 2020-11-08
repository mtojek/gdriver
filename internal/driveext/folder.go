package driveext

import (
	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"
)

func EnsureFolder(driveService *drive.Service, folderID string) error {
	file, err := driveService.Files.Get(folderID).Do()
	if err != nil {
		return errors.Wrapf(err, "can't read folder metadata (ID: %s)", folderID)
	}

	if file.MimeType != "application/vnd.google-apps.folder" {
		return errors.Wrapf(err, "resource is not a folder (ID: %s)", folderID)
	}
	return nil
}
