package main

import (
	"fmt"
	"regexp"
	"runtime"
	"unicode"
)

var getOpts = regexp.MustCompile(`^(--?)([^=]+)(.*?)$`)

func sortArgs(argsIn []string) (longArgs []string, args string, fileName string, path string, execArgs []string, err string) {

	var execArgsStorage []string
	if len(argsIn) < 2 {
		return argsIn, "", "", "", execArgsStorage, "ffind: too few arguments"
	}

	for i, argsElement := range argsIn {

		if argsElement == "-exec" {
			execArgsStorage = argsIn[i:]
			break
		}
		opts := getOpts.FindStringSubmatch(argsElement)
		if len(opts) == 0 {
			if len(fileName) == 0 {
				fileName = argsElement
			} else {
				path = argsElement
			}
		} else {
			switch opts[1] {
			case "-":
				args = opts[2]
				if len(opts) == 4 {
					args += opts[3]
				}
			case "--":
				longArgs = append(longArgs, opts[2])
			default:
				return argsIn, "", "", "", execArgsStorage, "ffind: invalid option input format"
			}
		}
	}
	return longArgs, args, fileName, path, execArgsStorage, ""
}

func makeCommand(longArgs []string, args string, fileName string, path string, execArgs []string) (argsOut []string, err string) {
	DebugLogger.Println("makeCommand(", longArgs, ",", args, ",", fileName, ",", path, ",", execArgs, ")")

	argsOut = append(argsOut, path) /* Adding path */

	preName, err := getArgs(args)
	if err != "" {
		return argsOut, err
	}
	argsOut = append(argsOut, preName...) /* Adding arguments going before filename */

	name, err := getName(fileName)
	if err != "" {
		return argsOut, err
	}
	argsOut = append(argsOut, name)        /* Adding filename */
	argsOut = append(argsOut, execArgs...) /* Adding arguments related to -exec */

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
			if runtime.GOOS == "openbsd" {
				return argsOut, fmt.Sprintf("ffind: -r: regex for find not supported in openbsd", opts)
			}
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
			return argsOut, fmt.Sprintf("ffind: -%c: unknown option", opts)
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
		return args, "ffind: depth option present but not specified"
	}
	return num, ""
}
