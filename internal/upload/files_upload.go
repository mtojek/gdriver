package upload

import (
	"fmt"
	"os"
	"strings"

	"github.com/avast/retry-go"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/api/drive/v3"

	"github.com/mtojek/gdriver/internal/osext"
)

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
		if err == nil {
			fmt.Println(" File verification failed")
			return errors.WithMessage(err, "remote file verification failed")
		}
		return nil
	}, retry.Attempts(3))
}

func mkdirAll(driveService *drive.Service, file *osext.LocalFile, folderID string) (string, error) {
	return "TODO", nil // TODO
}

func verifyFileIntegrity(driveService *drive.Service, localFile *osext.LocalFile, folderID string) error {
	q := fmt.Sprintf("trashed = false and '%s' in parents", folderID)
	files, err := driveService.Files.List().
		Fields("files(id, name, size, md5Checksum, mimeType, trashed)").
		Q(q).
		Do()
	if err != nil {
		return errors.Wrap(err, "files.list call failed")
	}
	if len(files.Files) != 1 {
		return fmt.Errorf("expected single item, got: %d", len(files.Files))
	}

	remoteFile := files.Files[0]

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

func uploadFileData(driveService *drive.Service, file *osext.LocalFile, parentID string, bar *progressbar.ProgressBar) error {
	fd, err := os.Open(file.Path)
	if err != nil {
		return errors.Wrapf(err, "can't open the local file (path: %s)", file.Path)
	}
	defer fd.Close()

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
