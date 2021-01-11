package upload

import (
	"fmt"
	"strings"

	"github.com/avast/retry-go"
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

		// TODO Navigate to the target, create subdirectories
		// TODO Upload file data
		// TODO Compare MD5
		fmt.Println(bar)

		return nil
	}, retry.Attempts(3))
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
