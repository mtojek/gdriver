package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func Authenticate(newCredentialsFile string) error {
	if newCredentialsFile != "" {
		err := writeCredentialsFile(newCredentialsFile)
		if err != nil {
			return errors.Wrap(err, "importing credentials failed")
		}
	}

	config, err := clientConfig()
	if err != nil {
		return errors.Wrap(err, "reading client config failed")
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

	err = writeTokenFile(authorizationToken)
	if err != nil {
		return errors.Wrap(err, "writing token file failed")
	}

	fmt.Println("Done.")
	return nil
}

func Client() (*http.Client, error) {
	token, err := readTokenFile()
	if err != nil {
		return nil, errors.Wrap(err, "reading token file failed")
	}

	c, err := clientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "creating client config failed")
	}
	return c.Client(context.Background(), token), nil
}

func Verify() error {
	_, err := readTokenFile()
	if err != nil {
		return errors.Wrap(err, "user hasn't been authenticated")
	}
	return nil
}

func clientConfig() (*oauth2.Config, error) {
	credentialsBody, err := readCredentialsFile()
	if err != nil {
		return nil, errors.Wrap(err, "reading credentials file filed")
	}

	config, err := google.ConfigFromJSON(credentialsBody, drive.DriveReadonlyScope)
	if err != nil {
		return nil, errors.Wrap(err, "parsing config file failed")
	}
	return config, nil
}
