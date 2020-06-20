package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	args := os.Args[1:]

	if len(args) < 3 {
		printUsage()
	}

	var opts string
	if args[0][0] == '-' {
		switch args[0][1] {
			case 'd':
				opts = "d"
			case 'f':
				opts = "f"
			default:
				opts = "f"
		}
	}

	cmd := exec.Command("find", args[2], "-type", opts, "-name", args[1])
	stdout, outErr := cmd.Output()
	if outErr != nil {
		panic(outErr)
	}
	fmt.Print(string(stdout))
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "find [-OPTION] ['FILENAME'] [PATH]")
	os.Exit(1)
}
