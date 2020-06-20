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
	args = getArgs(args)

	cmd := exec.Command("find", args...)
	stdout, outErr := cmd.Output()
	if outErr != nil {
		panic(outErr)
	}
	fmt.Print(string(stdout))
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "find -OPTION 'FILENAME' PATH")
	os.Exit(1)
}

func getArgs(osArgs []string) []string {
	var opts []string

	opts = append(opts, osArgs[2])

	opts = append(opts, "-type")
	if osArgs[0][0] == '-' {
		switch osArgs[0][1] {
		case 'd':
			opts = append(opts, "d")
		case 'f':
			opts = append(opts, "f")
		default:
			opts = append(opts, "f")
		}
	} else {
		printUsage()
	}

	opts = append(opts, "-name")
	opts = append(opts, osArgs[1])

	return opts
}
