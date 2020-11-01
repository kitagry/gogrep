package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/koron-go/prefixw"
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
