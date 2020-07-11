package main

import (
	"fmt"
	"regexp"
)

var getOpts = regexp.MustCompile(`^(-)(.*?)$`)

func getArgs(argsIn []string) (argsOut []string, err string) {
	argsOut = append(argsOut, argsIn[len(argsIn)-1])

	caseInsen := false
	regex := false

	if len(argsIn) == 3 {
		optsString := getOpts.FindStringSubmatch(argsIn[0])
		if len(optsString) > 2 {
			for _, opts := range optsString[2] {
				switch opts {
				case 'd':
					argsOut = append(argsOut, "-type", "d")
				case 'f':
					argsOut = append(argsOut, "-type", "f")
				case 'i':
					caseInsen = true
				case 'R':
					regex = true
				default:
					return argsOut, fmt.Sprintf("ffind: unsupported option '%c'", opts)
				}
			}
		} else {
			return argsOut, "ffind: Options must be preceded by '-'"
		}
	}

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
	argsOut = append(argsOut, argsIn[len(argsIn)-2])

	return argsOut, ""
}
