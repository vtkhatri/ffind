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

	// making command
	commandArgs, err := makeCommand(os.Args[1:]);
	if err != "" {
		fmt.Println(err)
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
