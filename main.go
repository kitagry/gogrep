package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
	"golang.org/x/sync/errgroup"
)

var isTerminal = isatty.IsTerminal(os.Stdout.Fd())

type searchQuery struct {
	Query  string
	Regexp bool

	Colorized bool

	Invert bool

	MaxCount int
}

func (s searchQuery) parse(args []string) (searchQuery, []string, error) {
	copy := s
	copy.Colorized = isTerminal
	if s.Query != "" {
		copy.Regexp = true
		return copy, args, nil
	}

	if len(args) == 0 {
		return copy, nil, fmt.Errorf("query should be set")
	}

	copy.Query = args[0]
	return copy, args[1:], nil
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage of %s:
	%s [OPTION]... PATTERNS [FILE]...
`, os.Args[0], os.Args[0])
}

func parseQueries(args []string, exp *string) (query string, files []string, useRegexp bool, err error) {
	if exp != nil && *exp != "" {
		return *exp, args, true, nil
	}

	if len(args) == 0 {
		return "", nil, false, fmt.Errorf("query should be set")
	}

	return args[0], args[1:], false, nil
}

func run() int {
	flag.Usage = func() {
		usage()
		flag.PrintDefaults()
	}
	sq := searchQuery{}
	flag.StringVar(&sq.Query, "e", "", "use PATTERNS for matching")
	flag.BoolVar(&sq.Invert, "v", false, "select non-matching lines")
	flag.IntVar(&sq.MaxCount, "m", 0, "stop after NUM selected lines")
	flag.Parse()
	args := flag.Args()

	var files []string
	var err error
	sq, files, err = sq.parse(args)
	if err != nil {
		usage()
		fmt.Fprintf(os.Stderr, "Try '%s --help' for more information.", os.Args[0])
		return 1
	}

	if len(files) == 0 {
		ch := make(chan string)
		var eg errgroup.Group
		eg.Go(func() error {
			err := Search(ch, os.Stdin, sq)
			if err != nil {
				return err
			}
			return nil
		})

		w := bufio.NewWriter(os.Stdout)
		for s := range ch {
			_, err := w.WriteString(s + "\n")
			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
				return 1
			}
		}
		err := w.Flush()
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return 1
		}

		err = eg.Wait()
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return 1
		}
	} else {
		err := searchFiles(files, sq)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
	}
	return 0
}

func main() {
	code := run()
	os.Exit(code)
}
