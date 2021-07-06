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

	// sort arguments into argument type
	longArgs, args, filename, path, execArgs, err := sortArgs(os.Args[1:])
	if err != "" {
		fmt.Println("Error sorting arguments:", err)
		os.Exit(1)
	}

	// handling -- flags, these are modifiers on how ffind behaves, need to be handled first
	err = longArgFlags(longArgs)
	if err != "" {
		fmt.Println("Error parsing long arguments:", err)
		os.Exit(1)
	}

	// compiling the command to be passed to find
	commandArgs, err := makeCommand(longArgs, args, filename, path, execArgs)
	if err != "" {
		fmt.Println("Error making command:", err)
		os.Exit(1)
	}

	// executing the command
	out, execErr := exec.Command("find", commandArgs...).CombinedOutput()
	if execErr != nil {
		fmt.Println("Error executing command:", string(out))
		os.Exit(1)
	}

	// outputting to terminal
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
			return fmt.Sprintf("Unsupported flag --%s", longArg)
		}
	}
	return ""
}
