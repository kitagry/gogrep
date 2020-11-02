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
		return searchWithRegexp(ch, data, q)
	}
	search(ch, data, q)
	return nil
}

func search(ch chan<- string, data io.Reader, q searchQuery) {
	defer close(ch)
	sc := bufio.NewScanner(data)
	sc.Split(bufio.ScanLines)

	i := 0
	for sc.Scan() {
		t := sc.Text()
		ind := strings.Index(t, q.Query)
		if (ind != -1 && q.Invert) || (ind == -1 && !q.Invert) {
			continue
		}

		i++
		if ind != -1 && !q.Invert {
			if q.Colorized {
				ch <- t[:ind] + colorizedText(t[ind:ind+len(q.Query)]) + t[ind+len(q.Query):]
			} else {
				ch <- t
			}
		} else if ind == -1 && q.Invert {
			ch <- t
		}

		if q.MaxCount != 0 && q.MaxCount <= i {
			break
		}
	}
}

func searchWithRegexp(ch chan<- string, data io.Reader, q searchQuery) error {
	defer close(ch)
	sc := bufio.NewScanner(data)
	sc.Split(bufio.ScanLines)

	mu, err := regexp.Compile(q.Query)
	if err != nil {
		return err
	}

	i := 0
	for sc.Scan() {
		t := sc.Text()
		match := mu.MatchString(t)

		if (match && q.Invert) || (!match && !q.Invert) {
			continue
		}

		i++
		if match && !q.Invert {
			if q.Colorized {
				strs := mu.FindAllString(t, -1)
				for _, str := range strs {
					t = strings.Replace(t, str, colorizedText(str), 1)
				}
				ch <- t
			} else {
				ch <- t
			}
		} else if !match && q.Invert {
			ch <- t
		}

		if q.MaxCount != 0 && q.MaxCount <= i {
			break
		}
	}
	return nil
}
