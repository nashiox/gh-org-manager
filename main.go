package main

import (
	"fmt"
	"os"
)

var version string

func main() {
	if err := execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
