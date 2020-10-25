package auth

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func Verify() error {
	configDir, err := configurationDir()
	if err != nil {
		return errors.Wrap(err, "reading configuration directory failed")
	}

	_, err = os.Stat(filepath.Join(configDir, tokenFile))
	if os.IsNotExist(err) {
		return errors.New("user hasn't been authenticated")
	}
	return nil
}
