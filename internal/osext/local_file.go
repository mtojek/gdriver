package osext

import (
	"fmt"

	"github.com/dustin/go-humanize"
)

type LocalFile struct {
	Name string
	Path string
	Size int64
}

func (lf *LocalFile) String() string {
	return fmt.Sprintf("%s (%s)", lf.Path, humanize.Bytes(uint64(lf.Size)))
}

type LocalFiles []*LocalFile

func (files LocalFiles) String() []string {
	var labels []string
	for _, file := range files {
		labels = append(labels, file.String())
	}
	return labels
}