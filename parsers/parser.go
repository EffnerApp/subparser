package parsers

import "github.com/EffnerApp/subparser/model"

type Parser interface {
	Parse(content string) ([]*model.Plan, error)
}
