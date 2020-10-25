package auth

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/mtojek/gdriver/internal/configuration"
)

func readCredentialsFile() ([]byte, error) {
	configDir, err := configuration.Dir()
	if err != nil {
		return nil, errors.Wrap(err, "reading configuration directory failed")
	}

	credentialsPath := filepath.Join(configDir, credentialsFile)
	credentialsBody, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		return nil, errors.Wrap(err, "reading credentials file failed. Use \"auth --import-credentials\" option.")
	}
	return credentialsBody, nil
}

func writeCredentialsFile(sourceFile string) error {
	c, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return errors.Wrap(err, "reading new credentials file failed")
	}

	configDir, err := configuration.Dir()
	if err != nil {
		return errors.Wrap(err, "reading configuration directory failed")
	}

	credentialsPath := filepath.Join(configDir, credentialsFile)
	err = ioutil.WriteFile(credentialsPath, c, 0644)
	if err != nil {
		return errors.Wrap(err, "writing credentials file failed")
	}
	return nil
}
