package filter

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type Filter struct {
	// package path regexp to include.
	includes []*regexp.Regexp

	// package path regexp to exclude.
	excludes []*regexp.Regexp
}

func New(includesStr, excludesStr string) (*Filter, error) {
	if includesStr == "" {
		return nil, errors.New("includesStr unspecified")
	}

	includes, err := parseRegexps(includesStr)
	if err != nil {
		return nil, fmt.Errorf("includesStr parse error: %v", err)
	}

	if excludesStr == "" {
		return &Filter{
			includes: includes,
		}, nil
	}

	excludes, err := parseRegexps(excludesStr)
	if err != nil {
		return nil, fmt.Errorf("excludesStr parse error: %v", err)
	}

	return &Filter{
		includes: includes,
		excludes: excludes,
	}, nil
}

func parseRegexps(input string) ([]*regexp.Regexp, error) {
	values := strings.Split(input, ",")
	result := make([]*regexp.Regexp, len(values))
	for i, ptn := range values {
		v, err := regexp.Compile(ptn)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %q: %w", ptn, err)
		}
		result[i] = v
	}
	return result, nil
}

func (f *Filter) IsTarget(value string) bool {
	return slices.ContainsFunc(f.includes, func(ptn *regexp.Regexp) bool {
		return ptn.MatchString(value)
	}) && !slices.ContainsFunc(f.excludes, func(ptn *regexp.Regexp) bool {
		return ptn.MatchString(value)
	})
}
