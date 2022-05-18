package main

import (
	"fmt"
	"regexp"
	"os"
	"unicode"
)

// SortedArg for code legibility
type SortedArgs struct {
	longArgs []string
	shortArgs string
	fileName string
	path string
	execArgs []string
}

var getOpts = regexp.MustCompile(`^(--?)([^=]+)(.*?)$`)

func makeCommand(argsIn []string) (commandArgs []string, retErr string) {

	// sort arguments into argument type
	sortedArgs, err := sortArgs(argsIn)
	if err != "" {
		retErr = fmt.Sprintf("Error sorting arguments: %s", err)
		return commandArgs, retErr
	}

	// handling -- flags, these are modifiers on how ffind behaves, need to be handled first
	err = longArgFlags(sortedArgs.longArgs)
	if err != "" {
		retErr = fmt.Sprintf("Error parsing long arguments: %s", err)
		return commandArgs, retErr
	}

	// compiling the command arguments following find call
	commandArgs, err = makeCommandArgs(sortedArgs)
	if err != "" {
		retErr = fmt.Sprintf("Error making command arguments: %s", err)
		return commandArgs, retErr
	}

	return commandArgs, ""
}

func sortArgs(argsIn []string) (sortedArgs SortedArgs, err string) {

	if len(argsIn) < 1 {
		return sortedArgs, "Too few arguments"
	}

	for i, argsElement := range argsIn {

		if argsElement == "-exec" {
			sortedArgs.execArgs = argsIn[i:]
			break
		}
		opts := getOpts.FindStringSubmatch(argsElement)
		if len(opts) == 0 {
			if len(sortedArgs.fileName) == 0 {
				sortedArgs.fileName = argsElement
			} else {
				sortedArgs.path = argsElement
			}
		} else {
			switch opts[1] {
			case "-":
				sortedArgs.shortArgs = opts[2]
				if len(opts) == 4 {
					sortedArgs.shortArgs += opts[3]
				}
			case "--":
				sortedArgs.longArgs = append(sortedArgs.longArgs, opts[2])
			default:
				return sortedArgs, "Invalid option input format"
			}
		}
	}

	return sortedArgs, ""
}

func makeCommandArgs(args SortedArgs) (argsOut []string, err string) {
	DebugLogger.Printf("makeCommandArgs(longArgs=%s, shortArgs=%s, fileName=%s, path=%s, execArgs=%s)", args.longArgs, args.shortArgs, args.fileName, args.path, args.execArgs)

	argsOut = append(argsOut, args.path)       /* Adding path */

	preName, err := getArgs(args.shortArgs)
	if err != "" {
		return argsOut, err
	}
	argsOut = append(argsOut, preName...)      /* Adding arguments going before filename */

	name, err := getName(args.fileName)
	if err != "" {
		return argsOut, err
	}
	argsOut = append(argsOut, name)    /* Adding filename */
	argsOut = append(argsOut, args.execArgs...) /* Adding arguments related to -exec */

	DebugLogger.Println("makeCommandArgs->", argsOut)
	return argsOut, ""
}

func getArgs(argsIn string) (argsOut []string, err string) {

	caseInsen := false
	regex := false

optionParsing:
	for _, opts := range argsIn {
		switch opts {
		case 'd':
			argsOut = append(argsOut, "-type", "d")
		case 'f':
			argsOut = append(argsOut, "-type", "f")
		case 'i':
			caseInsen = true
		case 'r':
			regex = true
		case 'e':
			depth, err := getDepth(argsIn)
			if err != "" {
				return argsOut, err
			}
			argsOut = append(argsOut, "-maxdepth", depth) /* Adding maxdepth level (accepted before filename) */
		case '=':
			break optionParsing
		case 'h':
			printUsage()
			os.Exit(0)
		default:
			return argsOut, fmt.Sprintf("Unknown option -%c", opts)
		}
	}
	argsOut = append(argsOut, globType(caseInsen, regex)...)

	return argsOut, ""
}

func globType(caseInsen bool, regex bool) (argsOut []string) {

	if caseInsen {
		if regex {
			argsOut = append(argsOut, "-iregex")
		} else {
			argsOut = append(argsOut, "-iname")
		}
	} else {
		if regex {
			argsOut = append(argsOut, "-regex")
		} else {
			argsOut = append(argsOut, "-name")
		}
	}
	return argsOut
}

func getName(nameIn string) (nameOut string, err string) {

	nameOut = nameIn
	for i := 0; i < len(nameIn); i++ {
		if (nameOut[i] == ' ') && (i != 0) {
			nameOut = nameOut[0:i] + "\\" + nameOut[i:]
			i++
		}
	}

	return nameOut, ""
}

func getDepth(args string) (num string, err string) {

	for i, char := range args {
		if unicode.IsDigit(char) {
			num += string(args[i])
		}
	}
	if len(num) == 0 {
		return args, "Depth option present but not specified"
	}
	return num, ""
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

func printUsage() {
	fmt.Println("Usage: ffind [-fdri] [-e=maxdepth] [--debug --help] [expression] [path]")
}
