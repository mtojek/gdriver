package upload

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"

	"github.com/mtojek/gdriver/internal/driveext"
	"github.com/mtojek/gdriver/internal/osext"
)

type FilesOptions struct {
	FolderID  string
	SourceDir string

	SelectionMode bool
}

func Files(driveService *drive.Service, options FilesOptions) error {
	err := checkSourceDir(options.SourceDir)
	if err != nil {
		return err
	}

	// If a resource path is provided, check if it refers to a folder.
	if options.FolderID != "" {
		fmt.Printf("Read folder metadata for \"%s\"\n", options.FolderID)

		err := driveext.EnsureFolder(driveService, options.FolderID)
		if err != nil {
			return err
		}
	}

	fmt.Println("List available files")
	files, err := osext.ListFiles(options.SourceDir)
	if err != nil {
		return errors.Wrap(err, "listing files failed")
	}

	if len(files) == 0 {
		fmt.Println("Source local folder is empty")
		return nil
	}

	if options.SelectionMode {
		fmt.Println("Select files to upload")
		files, err = selectFiles(files)
		if err != nil {
			return errors.Wrap(err, "can't select files to upload")
		}
	}

	fmt.Println("Upload files")
	for _, file := range files {
		err := uploadFile(driveService, file, options.FolderID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", errors.Wrapf(err, "uploading file \"%s\" failed", file.Path))
		}
	}
	return nil
}

func checkSourceDir(sourceDir string) error {
	fi, err := os.Stat(sourceDir)
	if err != nil {
		return errors.Wrap(err, "stat source dir failed")
	}

	if !fi.IsDir() {
		return errors.New("source directory must be a folder")
	}
	return nil
}
