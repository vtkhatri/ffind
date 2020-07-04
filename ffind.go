package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {

	if len(os.Args[1:]) < 2 {
		fmt.Println("too few arguments")
		printUsage()
	}

	cmd := exec.Command("find", parseArgs()...)
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("find failed with error %s", err)
	}
	fmt.Printf("%s", cmdOutput)
}

func printUsage() {
	log.Fatal("usage: ffind [-OPTIONS] NAME PATH")
}

func parseArgs() []string {
	var cmd []string

	cmd = append(cmd, os.Args[len(os.Args[1:])])

	caseInsen := false
	if os.Args[1][0] == '-' {
		for _, opts := range os.Args[1][1:] {
			switch opts {
			case 'd':
				cmd = append(cmd, "-type d")
			case 'f':
				cmd = append(cmd, "-type f")
			case 'i':
				caseInsen = true
			default:
				log.Fatal("Only -d, -f & -i supported")
			}
		}
	}

	if caseInsen {
		cmd = append(cmd, "-iname")
	} else {
		cmd = append(cmd, "-name")
	}
	cmd = append(cmd, "'"+os.Args[len(os.Args[1:])-1]+"'")

	fmt.Println(cmd)

	return cmd
}
