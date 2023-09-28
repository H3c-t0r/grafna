package sortopts

import (
	"fmt"
	"sort"
	"strings"

	"github.com/grafana/grafana/pkg/services/search/model"
	"github.com/grafana/grafana/pkg/util/errutil"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	// SortOptionsByQueryParam is a map to translate the "sort" query param values to SortOption(s)
	SortOptionsByQueryParam = map[string]model.SortOption{
		"name-asc":         newSortOption("name", false, 0),
		"name-desc":        newSortOption("name", true, 0),
		"email-asc":        newSortOption("email", false, 1),
		"email-desc":       newSortOption("email", true, 1),
		"memberCount-asc":  newIntSortOption("member_count", false, 2),
		"memberCount-desc": newIntSortOption("member_count", true, 2),
	}

	ErrorUnknownSortingOption = errutil.BadRequest("unknown sorting option")
)

type Sorter struct {
	Field         string
	Descending    bool
	WithTableName bool
}

func (s Sorter) OrderBy() string {
	table := "team."
	if !s.WithTableName {
		table = ""
	}
	if s.Descending {
		return fmt.Sprintf("%v%v DESC", table, s.Field)
	}
	return fmt.Sprintf("%v%v ASC", table, s.Field)
}

func newSortOption(field string, desc bool, index int) model.SortOption {
	direction := "asc"
	description := ("A-Z")
	if desc {
		direction = "desc"
		description = ("Z-A")
	}
	return model.SortOption{
		Name:        fmt.Sprintf("%v-%v", field, direction),
		DisplayName: fmt.Sprintf("%v (%v)", cases.Title(language.Und).String(field), description),
		Description: fmt.Sprintf("Sort %v in an alphabetically %vending order", field, direction),
		Index:       index,
		Filter:      []model.SortOptionFilter{Sorter{Field: field, Descending: desc, WithTableName: true}},
	}
}

func newIntSortOption(field string, desc bool, index int) model.SortOption {
	direction := "asc"
	description := ("Fewest-Most")
	if desc {
		direction = "desc"
		description = ("Most-Fewest")
	}
	return model.SortOption{
		Name:        fmt.Sprintf("%v-%v", field, direction),
		DisplayName: fmt.Sprintf("%v (%v)", cases.Title(language.Und).String(field), description),
		Description: fmt.Sprintf("Sort %v in a numerically %vending order", field, direction),
		Index:       index,
		Filter:      []model.SortOptionFilter{Sorter{Field: field, Descending: desc, WithTableName: false}},
	}
}

// ParseSortQueryParam parses the "sort" query param and returns an ordered list of SortOption(s)
func ParseSortQueryParam(param string) ([]model.SortOption, error) {
	opts := []model.SortOption{}
	if param != "" {
		optsStr := strings.Split(param, ",")
		for i := range optsStr {
			if opt, ok := SortOptionsByQueryParam[optsStr[i]]; !ok {
				return nil, ErrorUnknownSortingOption.Errorf("%v option unknown", optsStr[i])
			} else {
				opts = append(opts, opt)
			}
		}
		sort.Slice(opts, func(i, j int) bool {
			return opts[i].Index < opts[j].Index || (opts[i].Index == opts[j].Index && opts[i].Name < opts[j].Name)
		})
	}
	return opts, nil
}
