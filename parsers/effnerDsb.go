package parsers

import (
	"errors"
	"github.com/EffnerApp/subparser/model"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
	"time"
)

var ErrElementNotFound = errors.New("element not found")

// EffnerDSBParser uses https://github.com/PuerkitoBio/goquery to parse HTML into a substitution model
type EffnerDSBParser struct {
}

// Parse Parsing based on Sebi's implementation in "effnerapp-push-v3"
// https://github.com/EffnerApp/effnerapp-push-v3/blob/master/src/tools/dsbmobile/index.ts#L150
func (parser *EffnerDSBParser) Parse(content string) ([]*model.Plan, error) {
	documents, err := parseDocuments(content)

	if err != nil {
		return nil, err
	}

	plans := make([]*model.Plan, len(documents))

	for index, docHtml := range documents {
		plan, err := parsePlan(docHtml)

		if err != nil {
			return nil, err
		}

		plans[index] = plan
	}

	return plans, nil
}

func parsePlan(content string) (*model.Plan, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(content))

	if err != nil {
		return nil, err
	}

	// load information about the model
	date := findDate(document)
	title := findTitle(document)
	createdAt, err := findCreatedAt(document)

	if err != nil {
		return nil, err
	}

	// load the absent classes
	absent, err := findAbsent(document)

	if err != nil {
		return nil, err
	}

	// load the infos
	infos, err := findInfos(document)
	if err != nil {
		return nil, err
	}

	// load the substitutions
	substitutions, err := findSubstitutions(document)

	if err != nil {
		return nil, err
	}

	return &model.Plan{
		Title:         title,
		Date:          date,
		CreatedAt:     createdAt,
		Absent:        absent,
		Infos:         infos,
		Substitutions: substitutions,
	}, nil
}

func findInfos(document *goquery.Document) ([]string, error) {
	table := document.Find("table.F")

	if table == nil {
		return []string{}, nil
	}

	table = table.First()

	if table == nil {
		return []string{}, nil
	}

	infos := make([]string, 0)

	space := regexp.MustCompile(`\s+`)

	// very ugly but works...
	table.Find("th.F").Each(func(_ int, th *goquery.Selection) {
		text := th.Text()
		text = strings.TrimPrefix(text, "\n")

		if text == "" {
			return
		}

		text = strings.Replace(text, "\t", " ", -1)

		// some strange character they for some reason contain. Never touch this.
		text = strings.Replace(text, " ", " ", -1)

		text = strings.TrimSpace(text)
		text = space.ReplaceAllString(text, " ")

		infos = append(infos, text)
	})
	return infos, nil
}

func findSubstitutions(document *goquery.Document) ([]model.Substitution, error) {
	table := document.Find("table.k")
	if table != nil {
		// make an array of substitutions
		substitutions := make([]model.Substitution, 0)

		// loop through all "tbody" elements and parse the entries
		table.Find("tbody.k").Each(func(_ int, tbody *goquery.Selection) {
			// a tbody might contain one or more "tr.k"-elements. Each tr-element represents one substitution
			// but only the first tr-element has the class th-element
			class := strings.TrimSpace(tbody.Find("th.k").Text())

			// now we can parse one or more substitutions from the tr-elements
			tbody.Find("tr.k").Each(func(_ int, tr *goquery.Selection) {
				substitution := model.Substitution{
					Class: class,
				}

				// loop td-elements and bind them onto their representing attribute of the substitution
				tr.Find("td").Each(func(index int, elem *goquery.Selection) {
					switch index {
					case 0: // teacher
						substitution.Teacher = strings.TrimSpace(elem.Text())
					case 1: // period
						substitution.Period = strings.TrimSpace(elem.Text())
					case 2: // substitute
						substitution.Substitute = strings.TrimSpace(elem.Text())
					case 3: // room
						substitution.Room = strings.TrimSpace(elem.Text())
					case 4: // Info
						substitution.Info = strings.TrimSpace(elem.Text())
					}
				})
				substitutions = append(substitutions, substitution)
			})
		})
		return substitutions, nil
	}
	return nil, ErrElementNotFound
}

func findAbsent(document *goquery.Document) ([]model.Absent, error) {
	table := document.Find("table.K")
	if table != nil {
		absents := make([]model.Absent, 0)

		table.Find("tr.K").Each(func(i int, s *goquery.Selection) {
			// parse an absent out of the tbody selection
			// TODO Error handling if elements not found?
			class := strings.TrimSpace(s.Find("th.K").Text())
			periods := strings.TrimSpace(s.Find("td").Text())

			absents = append(absents, model.Absent{
				Class:   class,
				Periods: periods,
			})
		})
		return absents, nil
	}
	return nil, ErrElementNotFound
}

func findCreatedAt(document *goquery.Document) (time.Time, error) {
	elem := document.Find("h4")
	if elem != nil {
		text := strings.TrimSpace(elem.Text())
		text = text[11:strings.Index(text, ")")]

		// parse the text into a date-time
		date, err := time.Parse("_2.1.2006 um 15:04 Uhr", text)

		if err != nil {
			return time.Now(), err
		}
		return date, nil
	}
	return time.Now(), ErrElementNotFound
}

func findDate(document *goquery.Document) string {
	elem := document.Find("a[name]")
	if elem != nil {
		return elem.AttrOr("name", "")
	}
	return ""
}

func findTitle(document *goquery.Document) string {
	elem := document.Find("h2")
	if elem != nil {
		return strings.TrimSpace(elem.Text())
	}
	return ""
}

func parseDocuments(content string) ([]string, error) {
	// try parsing the content as HTML
	document, err := goquery.NewDocumentFromReader(strings.NewReader(content))

	if err != nil {
		return nil, err
	}

	startElements := make([]string, 0)

	// try to find the DATE tags, which are basically the sign for the start of a new "model"
	document.Find("a[name]").Each(func(i int, s *goquery.Selection) {
		name, exists := s.Attr("name")
		if exists && name != "oben" {
			// TODO Error Handling here? Error is very unlikely
			outer, _ := goquery.OuterHtml(s)

			startElements = append(startElements, outer)
		}
	})

	docHtml, err := document.Html()

	if err != nil {
		return nil, err
	}

	docSections := make([]string, len(startElements))

	for index, elemHtml := range startElements {
		elementIndex := strings.Index(docHtml, elemHtml)

		if index == len(startElements)-1 {
			// last element, just use the remaining HTML
			docSections[index] = docHtml[elementIndex:]
		} else {
			// find the end of the slice by using the start-index of the next documentStart element
			endIndex := strings.Index(docHtml, startElements[index+1])

			// create the section using the slice from currentStartIndex->nextStartIndex
			docSections[index] = docHtml[elementIndex:endIndex]
		}
	}
	return docSections, nil
}
