package destination

import (
	"encoding/json"
	"github.com/EffnerApp/subparser/model"
	"os"
)

type FileDestination struct {
	Path string
}

func (dest *FileDestination) Write(plans []*model.Plan) error {
	plansJson, err := json.Marshal(plans)

	if err != nil {
		return err
	}

	// write output to file
	file, err := os.Create(dest.Path)

	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(plansJson)
	if err != nil {
		return err
	}
	return nil
}
