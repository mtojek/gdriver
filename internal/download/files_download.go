package download

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/avast/retry-go"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/api/drive/v3"

	"github.com/mtojek/gdriver/internal/driveext"
)

func downloadFile(driveService *drive.Service, file *driveext.DriveFile, outputDir string) error {
	return retry.Do(func() error {
		bar := progressbar.DefaultBytes(
			file.Size,
			renderBarDescription(file),
		)

		state, err := driveext.EvaluateFileState(file, filepath.Join(outputDir, file.Path))
		if err != nil {
			fmt.Println(" File state evaluation failed.")
			return errors.Wrap(err, "file state evaluation failed")
		}

		bar.Set64(state.Offset())
		if ok, _ := state.Valid(); ok {
			return nil // file is complete and valid
		}

		err = downloadFileData(driveService, *state, bar)
		if err != nil {
			fmt.Println(" File download failed.")
			return errors.Wrap(err, "downloading file data failed")
		}

		state, err = driveext.EvaluateFileState(file, filepath.Join(outputDir, file.Path))
		if err != nil {
			fmt.Println(" File state evaluation failed.")
			return errors.Wrap(err, "file state evaluation failed")
		}

		if ok, err := state.Valid(); !ok {
			fmt.Println(" File verification failed.")
			return errors.Wrap(err, "file verification failed")
		}

		return nil
	}, retry.Attempts(3))
}

func downloadFileData(driveService *drive.Service, state driveext.FileState, bar *progressbar.ProgressBar) error {
	filesGetCall := driveService.Files.Get(state.File.Id)
	filesGetCall.Header().Add("Range", fmt.Sprintf("bytes=%d-", state.Offset()))
	resp, err := filesGetCall.Download()
	if err != nil {
		return errors.Wrap(err, "opening file stream failed")
	}

	// Create base directories
	basePath := filepath.Dir(state.LocalPath)
	err = os.MkdirAll(basePath, 0755)
	if err != nil {
		return errors.Wrapf(err, "creating base directories failed (path: %s)", basePath)
	}

	// Open destination file
	var flags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	if state.Offset() == 0 {
		flags = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}

	f, err := os.OpenFile(state.LocalPath, flags, 0644)
	if err != nil {
		return errors.Wrap(err, "opening destination file failed")
	}
	defer f.Close()

	// Download file data
	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	if err != nil {
		return errors.Wrap(err, "copying buffers failed")
	}
	return nil
}

func renderBarDescription(file *driveext.DriveFile) string {
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
