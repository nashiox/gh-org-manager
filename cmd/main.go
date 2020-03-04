package main

import (
	"fmt"
	"os"

	ghom "github.com/nashiox/gh-org-manager/cmd/gh-org-manager"
)

var version string

func main() {
	rootCmd := ghom.GetRootCmd(version)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
