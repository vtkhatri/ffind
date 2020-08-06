package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {

	if len(os.Args[1:]) < 2 {
		fmt.Println("too few arguments")
		printUsage()
		os.Exit(1)
	}

	args, err := getArgs(os.Args[1:])
	if err != "" {
		fmt.Println(err)
		os.Exit(1)
	}

	out, execErr := exec.Command("find", args...).CombinedOutput()
	if execErr != nil {
		fmt.Println(string(out))
		os.Exit(1)
	}
	fmt.Printf("%s", out)
}

func printUsage() {
	fmt.Println("usage: ffind [-OPTIONS] NAME PATH")
}
