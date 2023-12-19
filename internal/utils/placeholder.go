package utils

import (
	"regexp"
	"strconv"
	"strings"
)

type Command struct {
	Command string
	Args    []string
}

func GenerateCommand(command string, args []string) Command {
	nargs := len(args)
	indices := make([]int, nargs)
	r := regexp.MustCompile("%-?[0-9]+%")
	command = r.ReplaceAllStringFunc(command, func(s string) string {
		if j, err := strconv.Atoi(strings.Trim(s, "%")); err == nil {
			if j < 0 {
				j = nargs + j + 1
			}
			if j <= nargs && nargs > 0 {
				indices[j-1]++
				return args[j-1]
			} else {
				return ""
			}
		}
		return s
	})
	newArgs := make([]string, 0)
	for i, v := range indices {
		if v == 0 {
			newArgs = append(newArgs, args[i])
		}
	}
	return Command{
		Command: command,
		Args:    newArgs,
	}
}
