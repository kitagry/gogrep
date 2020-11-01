package main

import (
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
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
}

func TestSearchWithRegexp(t *testing.T) {
	data := strings.NewReader(`
aaa
bbb
ccc
bbbaaaccc
`)
	query := "a*c"
	ch := make(chan string)
	go searchWithRegexp(ch, data, query, false)

	results := make([]string, 0)
	for res := range ch {
		results = append(results, res)
	}

	if len(results) != 2 {
		t.Errorf("results length expect %d, got %d", 2, len(results))
	}

	if results[0] != "ccc" {
		t.Errorf(`results[0] expected "%s", got "%s"`, "ccc", results[0])
	}

	if results[1] != "bbbaaaccc" {
		t.Errorf(`results[1] expected "%s", got "%s"`, "bbbaaaccc", results[1])
	}
}
