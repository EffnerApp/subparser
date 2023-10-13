package destination

import (
	"encoding/json"
	"fmt"
	"subparser/parsers"
)

type SysoutDestination struct{}

func (*SysoutDestination) Write(plans []*parsers.Plan) error {
	plansJson, err := json.Marshal(plans)

	if err != nil {
		return err
	}

	fmt.Println(string(plansJson))

	return nil
}
