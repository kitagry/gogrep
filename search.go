package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func searchFiles(filenames []string, query string) error {
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("%s is not found", filename)
		}
		defer f.Close()

		ch := make(chan string)
		go search(ch, f, query)

		w := bufio.NewWriter(os.Stdout)
		for s := range ch {
			// TODO: use bytes for performance.
			if len(filenames) > 1 {
				w.WriteString(filename + ": ")
			}
			w.WriteString(s + "\n")
		}
		w.Flush()
	}
	return nil
}

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
