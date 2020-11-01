package download

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"
)

func listFiles(driveService *drive.Service, folderID, path string) (driveFiles, error) {
	var files []*driveFile
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
			if aFile.MimeType == "application/vnd.google-apps.folder" {
				fs, err := listFiles(driveService, aFile.Id, filepath.Join(path, aFile.Name))
				if err != nil {
					return nil, errors.Wrapf(err, "listing child folder failed (folderID: %s)", aFile.Id)
				}
				files = append(files, fs...)
				continue
			}

			if strings.HasPrefix(aFile.MimeType, "application/vnd.google-apps.") {
				continue // skip Google Docs
			}

			files = append(files, &driveFile{
				File: aFile,
				Path: filepath.Join(path, aFile.Name),
			})
		}

		nextPageToken = fileList.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	return files, nil
}
