package main

import (
	"fmt"
	"regexp"
	"unicode"
)

var getOpts = regexp.MustCompile(`^(--?)([^=]+)(.*?)$`)

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

func makeCommand(args SortedArgs) (argsOut []string, err string) {
	DebugLogger.Printf("makeCommand(longArgs=%s, shortArgs=%s, fileName=%s, path=%s, execArgs=%s)", args.longArgs, args.shortArgs, args.fileName, args.path, args.execArgs)

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

	DebugLogger.Println("makeCommand->", argsOut)
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
