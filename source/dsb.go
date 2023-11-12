package source

import "github.com/EffnerApp/subparser/dsb"

// DSBSource TODO Configuration like "load all plans", filter plans etc
type DSBSource struct {
	User string
	Pass string
}

func (dsbSrc *DSBSource) Load() (string, error) {
	// load the data from DSB
	dsbInstance := dsb.NewDSB(dsbSrc.User, dsbSrc.Pass)

	err := dsbInstance.Login()

	if err != nil {
		return "", err
	}

	err = dsbInstance.LoadTimetables()

	if err != nil {
		return "", err
	}

	// get the document from the dsb instance
	document := dsbInstance.Documents[0].Children[0]
	content, err := document.Download()

	if err != nil {
		return "", err
	}

	return string(content), nil
}
