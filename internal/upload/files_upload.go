package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/avast/retry-go"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/api/drive/v3"

	"github.com/mtojek/gdriver/internal/osext"
)

var errResourceNotFound = errors.New("resource not found")

func uploadFile(driveService *drive.Service, file *osext.LocalFile, folderID string) error {
	return retry.Do(func() error {
		bar := progressbar.DefaultBytes(
			file.Size,
			renderBarDescription(file),
		)

		parentID, err := mkdirAll(driveService, file, folderID)
		if err != nil {
			fmt.Println(" Creating directory tree failed.")
			return errors.Wrap(err, "creating directory tree failed")
		}

		err = verifyFileIntegrity(driveService, file, parentID)
		if err == nil {
			bar.Finish()
			return nil
		}

		err = uploadFileData(driveService, file, parentID, bar)
		if err != nil {
			fmt.Println(" File upload failed.")
			return errors.Wrap(err, "uploading file data failed")
		}

		err = verifyFileIntegrity(driveService, file, parentID)
		if err != nil {
			fmt.Println(" File verification failed")
			return errors.Wrap(err, "remote file verification failed")
		}
		return nil
	}, retry.Attempts(3))
}

func mkdirAll(driveService *drive.Service, file *osext.LocalFile, folderID string) (string, error) {
	parentID := folderID
	basePath := filepath.Dir(file.Path)
	for {
		if basePath == "." {
			break
		}

		first := firstDir(basePath)
		remoteResource, err := getSingleFile(driveService, first, parentID)
		if err != nil && err != errResourceNotFound {
			return "", errors.Wrapf(err, "can't get file (name: %s, parentID: %s)", first, parentID)
		}
		if err == errResourceNotFound {
			remoteResource, err = createFolder(driveService, first, parentID)
			if err != nil {
				return "", errors.Wrapf(err, "can't create folder (name: %s, parentID: %s)", first, parentID)
			}
		}

		if remoteResource.MimeType != "application/vnd.google-apps.folder" {
			return "", fmt.Errorf("folder expected, but found file (ID: %s, name: %s)", remoteResource.Id,
				remoteResource.Name)
		}
		parentID = remoteResource.Id

		basePath, err = filepath.Rel(first, basePath)
		if err != nil {
			return "", errors.Wrap(err, "filepath.Rel failed")
		}
	}

	if parentID == "" {
		return "root", nil
	}
	return parentID, nil
}

func firstDir(path string) string {
	i := strings.Index(path, "/")
	if i == -1 {
		return path
	}
	return path[:i]
}

func verifyFileIntegrity(driveService *drive.Service, localFile *osext.LocalFile, parentID string) error {
	remoteFile, err := getSingleFile(driveService, localFile.Name, parentID)
	if err != nil {
		return errors.Wrapf(err, "can't get file (name: %s, parentID: %s)", localFile.Name, parentID)
	}

	if localFile.Size != remoteFile.Size {
		return fmt.Errorf("remote file has different size (expected: %d, actual: %d",
			localFile.Size, remoteFile.Size)
	}

	if localFile.Md5Checksum != remoteFile.Md5Checksum {
		return fmt.Errorf("remote file has different checksum (expected: %s, actual: %s)",
			localFile.Md5Checksum, remoteFile.Md5Checksum)
	}
	return nil
}

func getSingleFile(driveService *drive.Service, name, parentID string) (*drive.File, error) {
	files, err := driveService.Files.List().
		Fields("files(id, name, size, md5Checksum, mimeType, trashed)").
		Q(fmt.Sprintf("trashed = false and '%s' in parents and name = '%s'", parentID, name)).
		Do()
	if err != nil {
		return nil, errors.Wrap(err, "files.list call failed")
	}
	if len(files.Files) == 0 {
		return nil, errResourceNotFound
	}
	return files.Files[0], nil
}

func createFolder(driveService *drive.Service, name, parentID string) (*drive.File, error) {
	resource, err := driveService.Files.
		Create(&drive.File{
			Name:     name,
			Parents:  []string{parentID},
			MimeType: "application/vnd.google-apps.folder",
		}).
		Fields("id", "name", "mimeType").
		Do()
	if err != nil {
		return nil, errors.Wrap(err, "files.create failed")
	}
	return resource, nil
}

func uploadFileData(driveService *drive.Service, file *osext.LocalFile, parentID string, bar *progressbar.ProgressBar) error {
	fd, err := os.Open(file.Path)
	if err != nil {
		return errors.Wrapf(err, "can't open the local file (path: %s)", file.Path)
	}
	defer fd.Close()

	bar.Set64(0)
	_, err = driveService.Files.
		Create(&drive.File{
			Name:    file.Name,
			Parents: []string{parentID},
		}).
		Media(fd).
		SupportsAllDrives(true).
		ProgressUpdater(func(current, total int64) {
			bar.Set64(current)
		}).
		Do()
	if err != nil {
		return errors.Wrap(err, "files.create failed")
	}
	bar.Finish()
	return nil
}

func renderBarDescription(file *osext.LocalFile) string {
	var description strings.Builder
	if len(file.Name) > 22 {
		description.WriteString(file.Name[:22])
		description.WriteString("...")
	} else {
		description.WriteString(fmt.Sprintf("%-25s", file.Name))
	}

	description.WriteString(" (MD5: ")
	description.WriteString(file.Md5Checksum)
	description.WriteByte(')')
	return description.String()
}
