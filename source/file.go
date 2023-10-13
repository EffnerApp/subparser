package source

import "os"

type FileSource struct {
	Path string
}

func (file *FileSource) Load() (string, error) {
	// load data from file
	content, err := os.ReadFile(file.Path)

	if err != nil {
		return "", err
	}

	return string(content), nil
}
