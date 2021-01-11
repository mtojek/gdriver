package upload

import "google.golang.org/api/drive/v3"

type FilesOptions struct {
	FolderID  string
	SourceDir string

	SelectionMode bool
}

func Files(driveService *drive.Service, options FilesOptions) error {
	return nil
}
