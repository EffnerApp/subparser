package destination

import "subparser/parsers"

type Destination interface {
	Write(plans []*parsers.Plan) error
}
