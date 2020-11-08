package driveext

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"
)

func ListFiles(driveService *drive.Service, folderID string) (DriveFiles, error) {
	return listFilesWithPath(driveService, folderID, "/")
}

func listFilesWithPath(driveService *drive.Service, folderID, path string) (DriveFiles, error) {
	var files []*DriveFile
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

		for _, file := range fileList.Files {
			if file.MimeType == "application/vnd.google-apps.folder" {
				fs, err := listFilesWithPath(driveService, file.Id, filepath.Join(path, file.Name))
				if err != nil {
					return nil, errors.Wrapf(err, "listing child folder failed (folderID: %s)", file.Id)
				}
				files = append(files, fs...)
				continue
			}

			if strings.HasPrefix(file.MimeType, "application/vnd.google-apps.") {
				continue // skip Google Docs
			}

			files = append(files, &DriveFile{
				File: file,
				Path: filepath.Join(path, file.Name),
			})
		}

		nextPageToken = fileList.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	return files, nil
}
