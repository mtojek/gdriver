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

		err := driveext.EnsureFolder(driveService, options.FolderID)
		if err != nil {
			return err
		}
	}

	fmt.Println("List available files")
	files, err := driveext.ListFiles(driveService, options.FolderID)
	if err != nil {
		return errors.Wrap(err, "listing files failed")
	}

	if len(files) == 0 {
		fmt.Println("Source Google Drive folder is empty")
		return nil
	}

	if options.SelectionMode {
		fmt.Println("Select files to download")
		files, err = selectFiles(files)
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
