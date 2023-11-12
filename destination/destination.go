package destination

import (
	"github.com/EffnerApp/subparser/model"
)

type Destination interface {
	Write(plans []*model.Plan) error
}
