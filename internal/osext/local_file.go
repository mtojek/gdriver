package osext

type LocalFile struct {
	Path string
}

func (lf *LocalFile) String() string {
	return "TODO"
}

type LocalFiles []*LocalFile

func (files LocalFiles) String() []string {
	var labels []string
	for _, file := range files {
		labels = append(labels, file.String())
	}
	return labels
}