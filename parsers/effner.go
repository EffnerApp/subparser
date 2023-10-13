package parsers

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"time"
)

var ErrElementNotFound = errors.New("element not found")

// EffnerParser uses https://github.com/PuerkitoBio/goquery to parse HTML into a substitution plan
type EffnerParser struct {
}

// Parse Parsing based on Sebi's implementation in "effnerapp-push-v3"
// https://github.com/EffnerApp/effnerapp-push-v3/blob/master/src/tools/dsbmobile/index.ts#L150
func (parser *EffnerParser) Parse(content string) ([]*Plan, error) {
	documents, err := parseDocuments(content)

	if err != nil {
		return nil, err
	}

	plans := make([]*Plan, len(documents))

	for index, docHtml := range documents {
		plan, err := parsePlan(docHtml)

		if err != nil {
			return nil, err
		}

		plans[index] = plan
	}

	return plans, nil
}

func parsePlan(content string) (*Plan, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(content))

	if err != nil {
		return nil, err
	}

	// load information about the plan
	date := findDate(document)
	title := findTitle(document)
	createdAt, _ := findCreatedAt(document) // TODO Handle error if created-at parsing failed? TIME INVALID

	// load the absent classes
	absent, err := findAbsent(document)

	if err != nil {
		return nil, err
	}

	// load the substitutions
	substitutions, err := findSubstitutions(document)

	if err != nil {
		return nil, err
	}

	return &Plan{
		Title:         title,
		Date:          date,
		CreatedAt:     createdAt,
		Absent:        absent,
		Substitutions: substitutions,
	}, nil
}

func findSubstitutions(document *goquery.Document) ([]Substitution, error) {
	table := document.Find("table.k")
	if table != nil {
		// make an array of substitutions
		substitutions := make([]Substitution, 0)

		// loop through all "tbody" elements and parse the entries
		table.Find("tbody.k").Each(func(_ int, tbody *goquery.Selection) {
			// a tbody might contain one or more "tr.k"-elements. Each tr-element represents one substitution
			// but only the first tr-element has the class th-element
			class := strings.TrimSpace(tbody.Find("th.k").Text())

			// now we can parse one or more substitutions from the tr-elements
			tbody.Find("tr.k").Each(func(_ int, tr *goquery.Selection) {
				substitution := Substitution{
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

func findAbsent(document *goquery.Document) ([]Absent, error) {
	table := document.Find("table.K")
	if table != nil {
		absents := make([]Absent, 0)

		table.Find("tr.K").Each(func(i int, s *goquery.Selection) {
			// parse an absent out of the tbody selection
			// TODO Error handling if elements not found?
			class := strings.TrimSpace(s.Find("th.K").Text())
			periods := strings.TrimSpace(s.Find("td").Text())

			absents = append(absents, Absent{
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
		date, err := time.Parse("02.01.2006 um 15:04 Uhr", text)

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

	// try to find the DATE tags, which are basically the sign for the start of a new "plan"
	document.Find("a[name]").Each(func(i int, s *goquery.Selection) {
		name, exists := s.Attr("name")
		if exists && name != "oben" {
			// TODO Error Handling here? Error is very unlikely
			outer, _ := goquery.OuterHtml(s)

			startElements = append(startElements, outer)
		}
	})

	documentHtml, err := document.Html()

	if err != nil {
		return nil, err
	}

	documentSections := make([]string, len(startElements))

	for index, elemHtml := range startElements {
		elementIndex := strings.Index(documentHtml, elemHtml)

		if index == len(startElements)-1 {
			// last element, just use the remaining HTML
			documentSections[index] = documentHtml[elementIndex:]
		} else {
			// find the end of the slice by using the start-index of the next documentStart element
			endIndex := strings.Index(documentHtml, startElements[index+1])

			// create the section using the slice from currentStartIndex->nextStartIndex
			documentSections[index] = documentHtml[elementIndex:endIndex]
		}
	}
	return documentSections, nil
}
