package cmd

import (
	"fmt"
	"os"
)

func Execute() {
	args := os.Args[1:]

	switch args[0] {
	case "ls":
		ls(args[1:])

	case "cat":
		cat(args[1:])

	default:
		fmt.Fprintf(os.Stderr, "unsupported command %s", args[0])
	}

}

func ls(params []string) {
	if len(params) < 1 {
		fmt.Fprintf(os.Stderr, "missing path to directory")
		return
	}

	dir, err := os.ReadDir(params[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read directory: %v", err)
		return
	}

	for _, v := range dir {
		fmt.Printf("%s\t", v.Name())
	}
}

func cat(params []string) {
	if len(params) < 1 {
		fmt.Fprintf(os.Stderr, "missing path to file")
		return
	}

	file, err := os.ReadFile(params[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read file: %v", err)
		return
	}

	fmt.Println(string(file))
}
