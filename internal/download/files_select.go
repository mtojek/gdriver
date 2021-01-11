package download

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"

	"github.com/mtojek/gdriver/internal/driveext"
)

func selectFiles(files driveext.DriveFiles) (driveext.DriveFiles, error) {
	fileSelectPrompt := &survey.MultiSelect{
		Message:  "Which files would you like to download",
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
