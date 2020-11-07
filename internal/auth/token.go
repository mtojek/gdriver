package auth

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"github.com/mtojek/gdriver/internal/configuration"
)

const (
	credentialsFile = "credentials.json"
	tokenFile       = "token.json"
)

func readTokenFile() (*oauth2.Token, error) {
	configDir, err := configuration.Dir()
	if err != nil {
		return nil, errors.Wrap(err, "reading configuration directory failed")
	}
	tokenPath := filepath.Join(configDir, tokenFile)
	tokenBody, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return nil, errors.Wrap(err, "reading token file failed")
	}

	var token oauth2.Token
	err = json.Unmarshal(tokenBody, &token)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshaling token failed")
	}
	return &token, nil
}

func writeTokenFile(authorizationToken *oauth2.Token) error {
	tokenBody, err := json.Marshal(authorizationToken)
	if err != nil {
		return errors.Wrap(err, "marshaling authorization token failed")
	}

	configDir, err := configuration.Dir()
	if err != nil {
		return errors.Wrap(err, "reading configuration directory failed")
	}

	err = ioutil.WriteFile(filepath.Join(configDir, tokenFile), tokenBody, 0644)
	if err != nil {
		return errors.Wrap(err, "writing token file failed")
	}
	return nil
}
