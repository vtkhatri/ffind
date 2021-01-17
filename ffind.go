package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
)

// DebugLogger is the global logging tool
var DebugLogger *log.Logger

func init() {
	DebugLogger = log.New(ioutil.Discard, "debug : ", log.Ldate|log.Ltime)
}

func main() {

	if runtime.GOOS != "openbsd" && runtime.GOOS != "linux" {
		fmt.Println("ffind is not supported in ", runtime.GOOS)
		os.Exit(0)
	}

	longArgs, args, filename, path, execArgs, err := sortArgs(os.Args[1:])
	if err != "" {
		fmt.Println(err)
		printUsage()
		os.Exit(1)
	}

	// for --debug and --help flags
	err = longArgFlags(longArgs)
	if err != "" {
		fmt.Println(err)
		printUsage()
		os.Exit(1)
	}

	commandArgs, err := makeCommand(longArgs, args, filename, path, execArgs)
	if err != "" {
		fmt.Println(err)
		printUsage()
		os.Exit(1)
	}

	out, execErr := exec.Command("find", commandArgs...).CombinedOutput()
	if execErr != nil {
		fmt.Println(string(out))
		printUsage()
		os.Exit(1)
	}
	fmt.Printf("%s", out)
}

func printUsage() {
	fmt.Println("Usage: ffind [-fdri] [-e=maxdepth] [--debug --help] [expression] [path]")
}

func longArgFlags(longArgs []string) string {
	for _, longArg := range longArgs {
		switch longArg {
		case "debug":
			DebugLogger.SetOutput(os.Stderr)
		case "help":
			printUsage()
			os.Exit(0)
		default:
			return fmt.Sprintf("ffind: unsupported flag --%s", longArg)
		}
	}
	return ""
}
