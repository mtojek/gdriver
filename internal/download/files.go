package download

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"

	"github.com/mtojek/gdriver/internal/driveext"
)

type FilesOptions struct {
	FolderID  string
	OutputDir string

	SelectionMode bool
}

func Files(driveService *drive.Service, options FilesOptions) error {
	err := checkOutputDir(options.OutputDir)
	if err != nil {
		return err
	}

	// If a resource path is provided, check if it refers to a folder.
	if options.FolderID != "" {
		fmt.Printf("Read folder metadata for \"%s\"\n", options.FolderID)

		file, err := driveService.Files.Get(options.FolderID).Do()
		if err != nil {
			return errors.Wrapf(err, "can't read folder metadata (ID: %s)", options.FolderID)
		}

		if file.MimeType != "application/vnd.google-apps.folder" {
			return errors.Wrapf(err, "resource is not a folder (ID: %s)", options.FolderID)
		}
	}

	fmt.Println("List available files")
	files, err := driveext.ListFiles(driveService, options.FolderID, "/")
	if err != nil {
		return errors.Wrap(err, "listing files failed")
	}

	if options.SelectionMode {
		fmt.Println("Select files to download")
		files, err = selectFilesToDownload(files)
		if err != nil {
			return errors.Wrap(err, "can't select files to download")
		}
	}

	fmt.Println("Download files")
	for _, file := range files {
		err := downloadFile(driveService, file, options.OutputDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", errors.Wrapf(err, "downloading file \"%s\" failed", file.Path))
		}
	}
	return nil
}

func checkOutputDir(outputDir string) error {
	fi, err := os.Stat(outputDir)
	if err != nil {
		return errors.Wrap(err, "stat output dir failed")
	}

	if !fi.IsDir() {
		return errors.New("output directory must be a folder")
	}
	return nil
}
