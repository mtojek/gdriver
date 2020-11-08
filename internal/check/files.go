package check

import (
	"fmt"
	"github.com/mitchellh/colorstring"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"

	"github.com/mtojek/gdriver/internal/driveext"
)

type FilesOptions struct {
	FolderID  string
	TargetDir string
}

func Files(driveService *drive.Service, options FilesOptions) error {
	err := checkTargetDir(options.TargetDir)
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

	fmt.Println("Check local files")

	var failed bool
	for _, file := range files {
		fullPath := filepath.Join(options.TargetDir, file.Path)
		fmt.Print(fullPath)
		fmt.Print(" ... ")

		state, err := driveext.EvaluateFileState(file, fullPath)
		if err != nil {
			failed = true

			colorstring.Printf("[red]%v\n", err)
			continue
		}
		_, err = state.Valid()
		if err != nil {
			failed = true

			colorstring.Printf("[red]%v\n", err)
			continue
		}
		colorstring.Println("[green]OK")
	}

	if failed {
		return errors.New("verification failed")
	}
	return nil
}

func checkTargetDir(targetDir string) error {
	fi, err := os.Stat(targetDir)
	if err != nil {
		return errors.Wrap(err, "stat target dir failed")
	}

	if !fi.IsDir() {
		return errors.New("target directory must be a folder")
	}
	return nil
}
