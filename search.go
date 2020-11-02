package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/koron-go/prefixw"
	"golang.org/x/sync/errgroup"
)

func searchFiles(filenames []string, sq searchQuery) error {
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("%s is not found", filename)
		}
		defer f.Close()

		ch := make(chan string)
		var eg errgroup.Group
		eg.Go(func() error {
			err := Search(ch, f, sq)
			if err != nil {
				return err
			}
			return nil
		})

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

		err = eg.Wait()
		if err != nil {
			return err
		}
	}
	return nil
}

func colorizedText(s string) string {
	return "\033[33;100m" + s + "\033[0m"
}

func Search(ch chan<- string, data io.Reader, q searchQuery) error {
	if q.Regexp {
		return searchWithRegexp(ch, data, q.Query, q.Colorized, q.Invert)
	}
	search(ch, data, q.Query, q.Colorized, q.Invert)
	return nil
}

func search(ch chan<- string, data io.Reader, query string, colorized, invert bool) {
	defer close(ch)
	sc := bufio.NewScanner(data)
	sc.Split(bufio.ScanLines)

	for sc.Scan() {
		t := sc.Text()
		ind := strings.Index(t, query)
		if ind != -1 && !invert {
			if colorized {
				ch <- t[:ind] + colorizedText(t[ind:ind+len(query)]) + t[ind+len(query):]
			} else {
				ch <- t
			}
		} else if ind == -1 && invert {
			ch <- t
		}
	}
}

func searchWithRegexp(ch chan<- string, data io.Reader, query string, colorized bool, invert bool) error {
	defer close(ch)
	sc := bufio.NewScanner(data)
	sc.Split(bufio.ScanLines)

	mu, err := regexp.Compile(query)
	if err != nil {
		return err
	}

	for sc.Scan() {
		t := sc.Text()
		match := mu.MatchString(t)
		if match && !invert {
			if colorized {
				strs := mu.FindAllString(t, -1)
				for _, str := range strs {
					t = strings.Replace(t, str, colorizedText(str), 1)
				}
				ch <- t
			} else {
				ch <- t
			}
		} else if !match && invert {
			ch <- t
		}
	}
	return nil
}
