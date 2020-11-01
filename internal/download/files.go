package download

import (
	"context"
	"fmt"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"github.com/mtojek/gdriver/internal/auth"
)

type driveFile struct {
	Path string
	*drive.File
}

func (df *driveFile) String() string {
	return fmt.Sprintf("%s (%s)", df.Path, humanize.Bytes(uint64(df.Size)))
}

type driveFiles []*driveFile

func (files driveFiles) String() []string {
	var labels []string
	for _, aFile := range files {
		labels = append(labels, aFile.String())
	}
	return labels
}

type FilesOptions struct {
	FolderID  string
	OutputDir string

	SelectionMode bool
}

func Files(options FilesOptions) error {
	err := checkOutputDir(options.OutputDir)
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

	// If a resource path is provided, check if it refers to a folder.
	if options.FolderID != "" {
		aFile, err := driveService.Files.Get(options.FolderID).Do()
		if err != nil {
			return errors.Wrapf(err, "can't read folder data (ID: %s)", options.FolderID)
		}

		if aFile.MimeType != "application/vnd.google-apps.folder" {
			return errors.Wrapf(err, "resource is not a folder (ID: %s)", options.FolderID)
		}
	}

	files, err := listFiles(driveService, options.FolderID, "/")
	if err != nil {
		return errors.Wrap(err, "listing files failed")
	}

	if options.SelectionMode {
		files, err = selectFilesToDownload(files)
		if err != nil {
			return errors.Wrap(err, "can't select files to download")
		}
	}

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