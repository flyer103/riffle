package main

import (
	"fmt"
	"os"

	"github.com/flyer103/riffle/pkg/riffle"
)

func main() {
	if err := riffle.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
