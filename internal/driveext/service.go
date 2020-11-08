package driveext

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"github.com/mtojek/gdriver/internal/auth"
)

func NewService() (*drive.Service, error) {
	oauthClient, err := auth.Client()
	if err != nil {
		return nil, errors.Wrap(err, "creating auth client failed")
	}

	driveService, err := drive.NewService(context.Background(), option.WithHTTPClient(oauthClient))
	if err != nil {
		return nil, errors.Wrap(err, "creating drive service failed")
	}
	return driveService, nil
}
