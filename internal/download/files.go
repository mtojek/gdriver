package download

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"github.com/mtojek/gdriver/internal/auth"
)

type structuredFile struct {
	Path string
	*drive.File
}

type structuredFiles []*structuredFile

func (files structuredFiles) String() []string {
	var labels []string
	for _, aFile := range files {
		labels = append(labels, fmt.Sprintf("%s (%s)", aFile.Path, humanize.Bytes(uint64(aFile.Size))))
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

	// List files in the folder.
	files, err := listFiles(driveService, options.FolderID, "/")
	if err != nil {
		return errors.Wrap(err, "listing files failed")
	}

	if options.SelectionMode {
		// Select files to download.
		fileSelectPrompt := &survey.MultiSelect{
			Message:  "Which files would you like to download?",
			Options:  files.String(),
			PageSize: 100,
		}

		var selected []string
		err = survey.AskOne(fileSelectPrompt, &selected, survey.WithValidator(survey.Required))
		if err != nil {
			return errors.Wrap(err, "file selection prompt failed")
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

func listFiles(driveService *drive.Service, folderID, path string) (structuredFiles, error) {
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
