package auth

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func configurationDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "reading user home directory failed")
	}

	configDir := filepath.Join(homeDir, gdriverFolder)
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return "", errors.Wrap(err, "creating gdriver configuration directory failed")
	}
	return configDir, nil
}
