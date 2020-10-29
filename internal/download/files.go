package download

import (
	"context"
	"fmt"
	"github.com/mtojek/gdriver/internal/auth"
	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"os"
	"path/filepath"
)

type structuredFile struct {
	Path string
	*drive.File
}

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

	// If a resource path is provided, check if it refers to a folder.
	if folderID != "" {
		aFile, err := driveService.Files.Get(folderID).Do()
		if err != nil {
			return errors.Wrapf(err, "can't read folder data (ID: %s)", folderID)
		}

		if aFile.MimeType != "application/vnd.google-apps.folder" {
			return errors.Wrapf(err, "resource is not a folder (ID: %s)", folderID)
		}
	}

	// List files in the folder.
	files, err := listFiles(driveService, folderID, "/")
	if err != nil {
		return errors.Wrap(err, "listing files failed")
	}

	for _, aFile := range files {
		fmt.Println(aFile.Path)
	}

	// TODO Selection mode
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

func listFiles(driveService *drive.Service, folderID, path string) ([]*structuredFile, error) {
	var files []*structuredFile
	var nextPageToken string
	for {
		q := "trashed = false"
		if folderID != "" {
			q = q + fmt.Sprintf(" and '%s' in parents", folderID)
		}

		filesListCall := driveService.Files.List().
			PageSize(100).
			Fields("nextPageToken, files(id, name, size, md5Checksum, mimeType, trashed)").
			OrderBy("name").
			Q(q)
		if nextPageToken != "" {
			filesListCall.PageToken(nextPageToken)
		}

		fileList, err := filesListCall.Do()
		if err != nil {
			return nil, errors.Wrap(err, "files.list call failed")
		}

		for _, aFile := range fileList.Files {
			if aFile.MimeType != "application/vnd.google-apps.folder" {
				files = append(files, &structuredFile{
					File: aFile,
					Path: filepath.Join(path, aFile.Name),
				})
				continue
			}

			fs, err := listFiles(driveService, aFile.Id, filepath.Join(path, aFile.Name))
			if err != nil {
				return nil, errors.Wrapf(err, "listing child folder failed (folderID: %s)", aFile.Id)
			}

			files = append(files, fs...)
		}

		nextPageToken = fileList.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	return files, nil
}
