package upload

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/mtojek/gdriver/internal/osext"
	"github.com/pkg/errors"
)

func selectFiles(files osext.LocalFiles) (osext.LocalFiles, error) {
	fileSelectPrompt := &survey.MultiSelect{
		Message:  "Which files would you like to upload",
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

func filterSelectedFiles(files osext.LocalFiles, selected []string) osext.LocalFiles {
	var filtered []*osext.LocalFile
	for _, s := range selected {
		for _, file := range files {
			if file.String() == s {
				filtered = append(filtered, file)
			}
		}
	}
	return filtered
}
