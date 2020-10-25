package download

import (
	"os"
	
	"github.com/pkg/errors"
)

func Files(folderID, outputDir string) error {
	fi, err := os.Stat(outputDir)
	if err != nil {
		return errors.Wrap(err, "stat output dir failed")
	}

	if !fi.IsDir() {
		return errors.New("output directory must be a folder")
	}
	return nil
}
