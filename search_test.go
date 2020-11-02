package main

import (
	"reflect"
	"strings"
	"testing"
)

var data string = `aaa
bbb
ccc
bbbaaaccc
`

func TestSearch(t *testing.T) {
	tests := map[string]struct {
		query    searchQuery
		expected []string
	}{
		"Normal search": {
			query: searchQuery{
				Query: "aaa",
			},
			expected: []string{"aaa", "bbbaaaccc"},
		},
		"Invert search": {
			query: searchQuery{
				Query:  "aaa",
				Invert: true,
			},
			expected: []string{"bbb", "ccc"},
		},
		"Regexp search": {
			query: searchQuery{
				Query:  "a*c",
				Regexp: true,
			},
			expected: []string{"ccc", "bbbaaaccc"},
		},
		"Regexp search with invert": {
			query: searchQuery{
				Query:  "a*c",
				Regexp: true,
				Invert: true,
			},
			expected: []string{"aaa", "bbb"},
		},
	}

	for m, tt := range tests {
		t.Run(m, func(t *testing.T) {
			ch := make(chan string)
			go Search(ch, strings.NewReader(data), tt.query)

			results := make([]string, 0)
			for result := range ch {
				results = append(results, result)
			}

			if !reflect.DeepEqual(results, tt.expected) {
				t.Errorf("expected %+v, got %+v", tt.expected, results)
			}
		})
	}
}
