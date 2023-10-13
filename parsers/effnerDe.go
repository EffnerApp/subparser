package parsers

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

// EffnerDEParser uses https://github.com/PuerkitoBio/goquery to parse HTML into a substitution plan
type EffnerDEParser struct {
}

// Parse parsing the substitutions from effner.de
func (parser *EffnerDEParser) Parse(content string) ([]*Plan, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(content))

	if err != nil {
		return nil, err
	}

	plans := make([]*Plan, 0)

	document.Find("h3").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		plan := Plan{
			Title:         title,
			Substitutions: make([]Substitution, 0),
		}

		// find the next table
		table := document.AfterSelection(s).Find("table").First()

		table.Find("tr").Each(func(i int, tr *goquery.Selection) {
			// skip the first, that's the table header
			if i == 0 {
				return
			}

			substitution := Substitution{}
			tr.Find("td").Each(func(i int, elem *goquery.Selection) {
				switch i {
				case 0:
					substitution.Class = strings.TrimSpace(elem.Text())
				case 1: // teacher
					substitution.Teacher = strings.TrimSpace(elem.Text())
				case 2: // period
					substitution.Period = strings.TrimSpace(elem.Text())
				case 3: // substitute
					substitution.Substitute = strings.TrimSpace(elem.Text())
				case 4: // room
					substitution.Room = strings.TrimSpace(elem.Text())
				case 5: // Info
					substitution.Info = strings.TrimSpace(elem.Text())
				}
			})

			plan.Substitutions = append(plan.Substitutions, substitution)
		})

		plans = append(plans, &plan)
	})

	return plans, nil
}
