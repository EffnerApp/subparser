package destination

import (
	"subparser/model"
)

type Destination interface {
	Write(plans []*model.Plan) error
}
