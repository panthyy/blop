package main

import (
	"fmt"
	"os"

	"github.com/panthyy/blop/pkg/blop"
)

func main() {

	rootCmd := blop.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
