package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func Authenticate(newCredentialsFile string) error {
	configDir, err := configurationDir()
	if err != nil {
		return errors.Wrap(err, "reading configuration directory failed")
	}

	credentialsPath := filepath.Join(configDir, credentialsFile)
	if newCredentialsFile != "" {
		c, err := ioutil.ReadFile(newCredentialsFile)
		if err != nil {
			return errors.Wrap(err, "reading new credentials file failed")
		}

		err = ioutil.WriteFile(credentialsPath, c, 0644)
		if err != nil {
			return errors.Wrap(err, "writing credentials file failed")
		}
	}

	credentialsBody, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		return errors.Wrap(err, "reading credentials file failed. Use --import-credentials option.")
	}

	config, err := google.ConfigFromJSON(credentialsBody, drive.DriveReadonlyScope)
	if err != nil {
		return errors.Wrap(err, "parsing config file failed")
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Open the following link in your browser, authenticate account and copy the authorization code: ")
	fmt.Println(authURL)
	fmt.Print("\nEnter authorization code: ")

	var authorizationCode string
	if _, err := fmt.Scan(&authorizationCode); err != nil {
		return errors.Wrap(err, "reading authorization code failed")
	}

	authorizationToken, err := config.Exchange(context.TODO(), authorizationCode)
	if err != nil {
		return errors.Wrap(err, "exchanging authorization code for token failed")
	}

	tokenBody, err := json.Marshal(authorizationToken)
	if err != nil {
		return errors.Wrap(err, "marshaling authorization token failed")
	}

	err = ioutil.WriteFile(filepath.Join(configDir, tokenFile), tokenBody, 0644)
	if err != nil {
		return errors.Wrap(err, "writing token file failed")
	}

	fmt.Println("Done.")
	return nil
}
