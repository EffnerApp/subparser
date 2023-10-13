package parsers

import "subparser/model"

type Parser interface {
	Parse(content string) ([]*model.Plan, error)
}
