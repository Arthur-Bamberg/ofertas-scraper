package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: ofertas-scraper run")
		os.Exit(2)
	}
	switch os.Args[1] {
	case "run":
		fmt.Fprintln(os.Stderr, "run: daily job wiring not implemented yet (domain + ValidarExtracao ready)")
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		os.Exit(2)
	}
}
