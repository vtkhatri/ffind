package main

import (
	"fmt"
	"regexp"
)

var getOpts = regexp.MustCompile(`^(-)(.*?)$`)

func getArgs(argsIn []string) (argsOut []string, err string) {
	argsOut = append(argsOut, argsIn[len(argsIn)-1]) /* Adding path */
	preName, err := getPreName(argsIn)
	if err != "" {
		return argsOut, err
	}
	argsOut = append(argsOut, preName...)            /* Adding arguments going before filename */
	argsOut = append(argsOut, argsIn[len(argsIn)-2]) /* Adding filename */

	return argsOut, ""
}

func getPreName(argsIn []string) (argsOut []string, err string) {
	caseInsen := false
	regex := false

	if len(argsIn) > 2 {
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
				case 'r':
					regex = true
				case 'e':
					argsOut = append(argsOut, "-maxdepth", argsIn[len(argsIn)-3]) /* Adding maxdepth level (accepted before filename) */
				default:
					return argsOut, fmt.Sprintf("ffind: unsupported option '%c'", opts)
				}
			}
		} else {
			return argsOut, "ffind: Options must be preceded by '-'"
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
