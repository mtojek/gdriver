package osext

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func ListFiles(root string) (LocalFiles, error) {
	var files []*LocalFile
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		files = append(files, &LocalFile{
			Name: info.Name(),
			Path: path,
			Size: info.Size(),
		})
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "walking through files failed (root: %s)", root)
	}
	return files, nil
}
