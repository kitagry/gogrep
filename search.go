package main

import (
	"bufio"
	"io"
	"strings"
)

func search(ch chan<- string, data io.Reader, query string) {
	sc := bufio.NewScanner(data)
	sc.Split(bufio.ScanLines)

	for sc.Scan() {
		t := sc.Text()
		if strings.Contains(t, query) {
			ch <- t
		}
	}
	close(ch)
}
