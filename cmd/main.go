package main

import (
	"fmt"
	"os"

	"github.com/nashiox/gh-org-manager/cmd/ghom/cmd"
)

var version string

func main() {
	rootCmd := cmd.GetRootCmd(version)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
