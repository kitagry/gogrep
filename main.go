package main

import (
	"bufio"
	"fmt"
	"os"
)

func run() int {
	ch := make(chan string)
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "arguments should be set\n")
		return 1
	}

	query := os.Args[1]
	go search(ch, os.Stdin, query)

	w := bufio.NewWriter(os.Stdout)
	for s := range ch {
		w.WriteString(s + "\n")
	}
	w.Flush()
	return 0
}

func main() {
	code := run()
	os.Exit(code)
}
