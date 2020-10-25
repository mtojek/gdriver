package download

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"github.com/mtojek/gdriver/internal/auth"
)

func Files(folderID, outputDir string) error {
	err := checkOutputDir(outputDir)
	if err != nil {
		return err
	}

	oauthClient, err := auth.Client()
	if err != nil {
		return errors.Wrap(err, "creating auth client failed")
	}

	driveService, err := drive.NewService(context.Background(), option.WithHTTPClient(oauthClient))
	if err != nil {
		return errors.Wrap(err, "creating drive service failed")
	}
	fmt.Println(driveService)

	// TODO folderID as resourceID, selection mode
	// TODO List all files (next token) in the given folder (flat)
	// TODO For every file:
	// TODO 	Check if local file exists in output directory
	// TODO 	Check MD5 remote vs local file
	// TODO 	Check if size(local file) < size(remote local)
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
