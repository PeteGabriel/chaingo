package main

import (
	"fmt"
	"os"
)

func main() {

	if err := run(); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}

func run() error {
	cli := CLI{}
	cli.Run()
	return nil
}
