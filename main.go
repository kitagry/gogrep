package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage of %s:
	%s [OPTION]... PATTERNS [FILE]...
`, os.Args[0], os.Args[0])
}

func run() int {
	flag.Usage = func() {
		usage()
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()

	ch := make(chan string)
	if len(args) == 0 {
		usage()
		fmt.Fprintf(os.Stderr, "Try '%s --help' for more information.", os.Args[0])
		return 1
	}

	query := args[0]
	if len(args) == 1 {
		go search(ch, os.Stdin, query)

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
		err := searchFiles(args[1:], query)
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
