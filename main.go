package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func run() int {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
	%s [OPTION]... PATTERNS [FILE]...
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()

	ch := make(chan string)
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "arguments should be set\n")
		return 1
	}

	query := args[0]
	if len(args) == 1 {
		go search(ch, os.Stdin, query)

		w := bufio.NewWriter(os.Stdout)
		for s := range ch {
			w.WriteString(s + "\n")
		}
		w.Flush()
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
