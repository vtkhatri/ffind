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

	longArgs, args, filename, path, err := sortArgs(os.Args[1:])
	if err != "" {
		fmt.Println(err)
		os.Exit(1)
	}

	err = longArgFlags(longArgs)
	if err != "" {
		fmt.Println(err)
		os.Exit(1)
	}

	commandArgs, err := makeCommand(longArgs, args, filename, path)
	if err != "" {
		fmt.Println(err)
		os.Exit(1)
	}

	out, execErr := exec.Command("find", commandArgs...).CombinedOutput()
	if execErr != nil {
		fmt.Println(string(out))
		os.Exit(1)
	}
	fmt.Printf("%s", out)
}

func printUsage() {
	fmt.Println("usage: ffind [-OPTIONS] NAME PATH")
}

func longArgFlags(longArgs []string) string {
	for _, longArg := range longArgs {
		switch longArg {
		case "debug":
			/* to implement debug logic */
		case "help":
			printUsage()
			os.Exit(1)
		default:
			return fmt.Sprintf("unsupported flag --%s", longArg)
		}
	}
	return ""
}
