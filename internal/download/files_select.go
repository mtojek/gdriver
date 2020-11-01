package download

import (
	"github.com/pkg/errors"

	"github.com/AlecAivazis/survey/v2"
)

func selectFilesToDownload(files driveFiles) (driveFiles, error) {
	fileSelectPrompt := &survey.MultiSelect{
		Message:  "Which files would you like to download?",
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

func filterSelectedFiles(files driveFiles, selected []string) driveFiles {
	var filtered []*driveFile
	for _, s := range selected {
		for _, aFile := range files {
			if aFile.String() == s {
				filtered = append(filtered, aFile)
			}
		}
	}
	return filtered
}
