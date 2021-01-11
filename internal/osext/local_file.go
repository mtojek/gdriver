package osext

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
)

type LocalFile struct {
	Name string
	Path string
	Size int64
	Md5Checksum string
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

func (files LocalFiles) CalculateMd5Checksums() error {
	var err error
	for _, f := range files {
		f.Md5Checksum, err = Md5Checksum(f.Path)
		if err != nil {
			return errors.Wrapf(err, "calculating MD5 checksum failed (path: %s)", f.Path)
		}
	}
	return nil
}