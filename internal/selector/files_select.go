package selector

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"

	"github.com/mtojek/gdriver/internal/driveext"
)

func SelectFiles(files driveext.DriveFiles, action string) (driveext.DriveFiles, error) {
	fileSelectPrompt := &survey.MultiSelect{
		Message:  fmt.Sprintf("Which files would you like to %s?", action),
		Options:  files.String(),
		PageSize: 20,
	}

	var selected []string
	err := survey.AskOne(fileSelectPrompt, &selected, survey.WithValidator(survey.Required))
	if err != nil {
		return nil, errors.Wrap(err, "prompt failed")
	}

	files = filterSelectedFiles(files, selected)
	return files, nil
}

func filterSelectedFiles(files driveext.DriveFiles, selected []string) driveext.DriveFiles {
	var filtered []*driveext.DriveFile
	for _, s := range selected {
		for _, file := range files {
			if file.String() == s {
				filtered = append(filtered, file)
			}
		}
	}
	return filtered
}
