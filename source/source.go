package source

import "errors"

type Source interface {
	Load() (string, error)
}

func GetSourceFromType(sourceType string, user string, pass string, path string) (Source, error) {
	switch sourceType {
	case "dsb":
		return &DSBSource{
			User: user,
			Pass: pass,
		}, nil
	case "file":
		return &FileSource{
			Path: path,
		}, nil
	case "effner":
		return &EffnerDESource{
			Password: pass,
		}, nil
	default:
		return nil, errors.New("unknown source type '" + sourceType + "'")
	}
}
