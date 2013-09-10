package goschedule

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

// Position provides the start and end indices of a struct
// extracted from a document.
type position struct {
	Start int
	End   int
}

// A College is a UW college that has many departments
type College struct {
	Name         string
	Abbreviation string `pk:"true"`
	position     `ignore:"true"`
}

// A Department is a UW department that has many classes.
type Dept struct {
	CollegeKey   string `fk:"College"`
	Name         string
	Abbreviation string `pk:"true"`
	Link         string
}

// ExtractAbbreviation extracts a department abbreviation from content (assumed
// to be a class index of a department) and sets the Dept.Abbreviation attribute.
//
// Returns an error if no department abbreviation was found.
func (d *Dept) ScrapeAbbreviation(content string) error {
	if match := regexp.MustCompile(`(?i)<a name=.+?>(.+?)&nbsp;&nbsp; .*?</a>`).FindStringSubmatch(content); match == nil || len(match) < 2 {
		return fmt.Errorf("No match found in content")
	} else {
		abbreviation := html.UnescapeString(match[1])
		abbreviation = strings.ToLower(strings.Replace(abbreviation, " ", "", -1))
		d.Abbreviation = abbreviation
	}
	return nil
}

// A Class is UW class that has many sections.
type Class struct {
	DeptKey          string `fk:"Dept"`
	AbbreviationCode string `pk:"true"`
	Abbreviation     string
	Code             string
	Title            string
	Description      string
	position         `ignore:"true"`
}

// A Sect is a UW section.
type Sect struct {
	ClassKey     string `fk:"Class"`
	Restriction  string
	SLN          string `pk:"true"`
	Section      string
	Credit       string
	MeetingTimes string // JSON representation
	Instructor   string
	Status       string
	TakenSpots   int64
	TotalSpots   int64
	Grades       string
	Fee          string
	Other        string
	Info         string
}

// A MeetingTime represents when a Sect is held. Some Sect's
// have multiple meeting times.
type MeetingTime struct {
	Days     string
	Time     string
	Building string
	Room     string
}