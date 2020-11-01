package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/koron-go/prefixw"
	"github.com/mattn/go-isatty"
)

func searchFiles(filenames []string, query string) error {
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("%s is not found", filename)
		}
		defer f.Close()

		ch := make(chan string)
		go search(ch, f, query, isatty.IsTerminal(os.Stdout.Fd()))

		var w io.Writer
		bw := bufio.NewWriter(os.Stdout)
		w = bw
		if len(filenames) > 1 {
			w = prefixw.New(w, filename+": ")
		}
		for s := range ch {
			_, err = w.Write([]byte(s + "\n"))
			if err != nil {
				return err
			}
		}
		err = bw.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}

func search(ch chan<- string, data io.Reader, query string, colorized bool) {
	sc := bufio.NewScanner(data)
	sc.Split(bufio.ScanLines)

	for sc.Scan() {
		t := sc.Text()
		if ind := strings.Index(t, query); ind != -1 {
			if colorized {
				ch <- t[:ind] + "\033[33;100m" + t[ind:ind+len(query)] + "\033[0m" + t[ind+len(query):]
			} else {
				ch <- t
			}
		}
	}
	close(ch)
}
