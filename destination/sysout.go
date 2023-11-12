package destination

import (
	"encoding/json"
	"fmt"
	"github.com/EffnerApp/subparser/model"
)

type SysoutDestination struct{}

func (*SysoutDestination) Write(plans []*model.Plan) error {
	plansJson, err := json.Marshal(plans)

	if err != nil {
		return err
	}

	fmt.Println(string(plansJson))

	return nil
}
