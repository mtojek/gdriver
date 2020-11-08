package driveext

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"google.golang.org/api/drive/v3"
)

type DriveFile struct {
	Path string
	*drive.File
}

func (df *DriveFile) String() string {
	return fmt.Sprintf("%s (%s)", df.Path, humanize.Bytes(uint64(df.Size)))
}

type DriveFiles []*DriveFile

func (files DriveFiles) String() []string {
	var labels []string
	for _, file := range files {
		labels = append(labels, file.String())
	}
	return labels
}
