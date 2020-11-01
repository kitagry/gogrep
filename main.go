package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
)

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
	exp := flag.String("e", "", "use PATTERNS for matching")
	flag.Parse()
	args := flag.Args()

	query, files, useRegexp, err := parseQueries(args, exp)
	if err != nil {
		usage()
		fmt.Fprintf(os.Stderr, "Try '%s --help' for more information.", os.Args[0])
		return 1
	}

	if len(files) == 0 {
		ch := make(chan string)
		if useRegexp {
			go searchWithRegexp(ch, os.Stdin, query, isatty.IsTerminal(os.Stdout.Fd()))
		} else {
			go search(ch, os.Stdin, query, isatty.IsTerminal(os.Stdout.Fd()))
		}

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
	} else {
		err := searchFiles(files, query, useRegexp)
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
