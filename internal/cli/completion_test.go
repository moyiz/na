package cli

import (
	"slices"
	"testing"

	"github.com/moyiz/na/internal/config"
)

type ListNextPartsTestCase struct {
	Aliases  []config.Alias
	Prefix   []string
	Expected []string
}

func TestListNextParts(t *testing.T) {
	var result []string

	aliases := []config.Alias{
		{
			Name:    "my command",
			Command: "cmd1",
		},
		{
			Name:    "my other command",
			Command: "cmd1",
		},
		{
			Name:    "my other command2",
			Command: "cmd2",
		},
		{
			Name:    "another command2",
			Command: "cmd1",
		},
	}
	cases := []ListNextPartsTestCase{
		{
			Aliases:  aliases,
			Prefix:   []string{},
			Expected: []string{"my", "another"},
		},
		{
			Aliases:  aliases,
			Prefix:   []string{"my"},
			Expected: []string{"command", "other"},
		},
		{
			Aliases:  aliases,
			Prefix:   []string{"my", "other"},
			Expected: []string{"command", "command2"},
		},
		{
			Aliases:  aliases,
			Prefix:   []string{"another", "command2"},
			Expected: []string{},
		},
		{
			Aliases:  aliases,
			Prefix:   []string{"non", "existing"},
			Expected: []string{},
		},
	}
	for _, testCase := range cases {
		result = ListNextParts(aliases, testCase.Prefix)
		if slices.Compare(result, testCase.Expected) != 0 {
			t.Errorf("Expected %v but got %v", testCase.Expected, result)
		}
	}
}
