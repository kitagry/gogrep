package main

import (
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	t.Run("Two match items", func(t *testing.T) {
		data := strings.NewReader(`
aaa
bbb
ccc
aaa
bbbaaaccc
`)
		query := "aaa"
		ch := make(chan string)
		go search(ch, data, query, false)

		results := make([]string, 0)
		for result := range ch {
			results = append(results, result)
		}

		if len(results) != 3 {
			t.Errorf("results length is wrong: %d", len(results))
		}

		if results[0] != "aaa" {
			t.Errorf(`result expected "%s", go "%s"`, query, results[0])
		}

		if results[1] != "aaa" {
			t.Errorf(`result expected "%s", go "%s"`, query, results[1])
		}

		if results[2] != "bbbaaaccc" {
			t.Errorf(`result expected "%s", go "%s"`, query, results[2])
		}
	})
}
