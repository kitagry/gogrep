package main

import (
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	data := `aaa
bbb
ccc
aaa
bbbaaaccc
`
	query := "aaa"
	t.Run("normal search", func(t *testing.T) {
		ch := make(chan string)
		go search(ch, strings.NewReader(data), query, false, false)

		results := make([]string, 0)
		for result := range ch {
			results = append(results, result)
		}

		if len(results) != 3 {
			t.Fatalf("results length is wrong: %d", len(results))
		}

		if results[0] != "aaa" {
			t.Errorf(`result expected "%s", go "%s"`, "aaa", results[0])
		}

		if results[1] != "aaa" {
			t.Errorf(`result expected "%s", go "%s"`, "aaa", results[1])
		}

		if results[2] != "bbbaaaccc" {
			t.Errorf(`result expected "%s", go "%s"`, "bbbaaaccc", results[2])
		}
	})

	t.Run("invert search", func(t *testing.T) {
		ch := make(chan string)
		go search(ch, strings.NewReader(data), query, false, true)

		results := make([]string, 0)
		for result := range ch {
			results = append(results, result)
		}

		if len(results) != 2 {
			t.Fatalf("results length is wrong: %d", len(results))
		}

		if results[0] != "bbb" {
			t.Errorf(`result expected "%s", go "%s"`, "bbb", results[0])
		}

		if results[1] != "ccc" {
			t.Errorf(`result expected "%s", go "%s"`, "ccc", results[1])
		}
	})
}

func TestSearchWithRegexp(t *testing.T) {
	data := `aaa
bbb
ccc
bbbaaaccc
`
	query := "a*c"
	t.Run("normal search", func(t *testing.T) {
		ch := make(chan string)
		go searchWithRegexp(ch, strings.NewReader(data), query, false, false)

		results := make([]string, 0)
		for res := range ch {
			results = append(results, res)
		}

		if len(results) != 2 {
			t.Fatalf("results length expect %d, got %d", 2, len(results))
		}

		if results[0] != "ccc" {
			t.Errorf(`results[0] expected "%s", got "%s"`, "ccc", results[0])
		}

		if results[1] != "bbbaaaccc" {
			t.Errorf(`results[1] expected "%s", got "%s"`, "bbbaaaccc", results[1])
		}
	})

	t.Run("invert search", func(t *testing.T) {
		ch := make(chan string)
		go searchWithRegexp(ch, strings.NewReader(data), query, false, true)

		results := make([]string, 0)
		for res := range ch {
			results = append(results, res)
		}

		if len(results) != 2 {
			t.Fatalf("results length expect %d, got %d", 2, len(results))
		}

		if results[0] != "aaa" {
			t.Errorf(`results[0] expected "%s", got "%s"`, "aaa", results[0])
		}

		if results[1] != "bbb" {
			t.Errorf(`results[1] expected "%s", got "%s"`, "bbb", results[1])
		}
	})
}
