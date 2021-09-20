package schema

import (
	"errors"
	"regexp"
)

const (
	ProjectStatusTodo    = "todo"
	ProjectStatusHacking = "hacking"
	ProjectStatusDone    = "done"

	BlacklistRegex = "regex"
	BlacklistValue = "value"
)

type Blacklist struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

func NewBlacklist(value string, blacklistType string) (*Blacklist, error) {
	if blacklistType != BlacklistRegex && blacklistType != BlacklistValue {
		return nil, errors.New("Blacklist type should be regex or value only")
	}
	return &Blacklist{
		Value: value,
		Type:  blacklistType,
	}, nil
}

func BlacklistMatch(origin string, blacklist []Blacklist) bool {
	for _, b := range blacklist {
		switch b.Type {
		case BlacklistRegex:
			r, err := regexp.Compile(b.Value)
			if err != nil {
				continue
			}
			if r.MatchString(origin) {
				return true
			}
		case BlacklistValue:
			if origin == b.Value {
				return true
			}
		}
	}
	return false
}

func FilterBlacklistMatch(origins []string, blacklist []Blacklist) []string {
	result := []string{}
	for _, origin := range origins {
		if !BlacklistMatch(origin, blacklist) {
			result = append(result, origin)
		}
	}
	return result
}
