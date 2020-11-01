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

func searchFiles(filenames []string, query string, useRegexp bool) error {
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("%s is not found", filename)
		}
		defer f.Close()

		ch := make(chan string)
		var eg errgroup.Group
		if useRegexp {
			eg.Go(func() error {
				err := searchWithRegexp(ch, f, query, isTerminal)
				if err != nil {
					return err
				}
				return nil
			})
		} else {
			go search(ch, f, query, isTerminal)
		}

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

func search(ch chan<- string, data io.Reader, query string, colorized bool) {
	defer close(ch)
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
}

func searchWithRegexp(ch chan<- string, data io.Reader, query string, colorized bool) error {
	defer close(ch)
	sc := bufio.NewScanner(data)
	sc.Split(bufio.ScanLines)

	mu, err := regexp.Compile(query)
	if err != nil {
		return err
	}

	for sc.Scan() {
		t := sc.Text()
		if mu.MatchString(t) {
			if colorized {
				strs := mu.FindAllString(t, -1)
				for _, str := range strs {
					t = strings.Replace(t, str, "\033[33;100m"+str+"\033[0m", 1)
				}
				ch <- t
			} else {
				ch <- t
			}
		}
	}
	return nil
}
