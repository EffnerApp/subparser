package parsers

import (
	"errors"
	"github.com/EffnerApp/subparser/model"
)

type Parser interface {
	Parse(content string) ([]*model.Plan, error)
}

func GetParserFromType(parserType string) (Parser, error) {
	switch parserType {
	case "effner":
		return &EffnerDSBParser{}, nil
	case "effner-de":
		return &EffnerDEParser{}, nil
	default:
		return nil, errors.New("unknown parser type '" + parserType + "'")
	}
}
