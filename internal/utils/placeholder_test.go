package utils

import (
	"slices"
	"testing"
)

type GenerateCommandTestCase struct {
	Command  string
	Args     []string
	Expected Command
}

func TestGenerateCommand(t *testing.T) {
	for _, testCase := range []GenerateCommandTestCase{
		{
			Command: "my command",
			Args:    []string{"some", "args"},
			Expected: Command{
				Command: "my command",
				Args:    []string{"some", "args"},
			},
		},
		{
			Command: "my %1%",
			Args:    []string{"some", "args"},
			Expected: Command{
				Command: "my some",
				Args:    []string{"args"},
			},
		},
		{
			Command: "my %2%",
			Args:    []string{"some", "args"},
			Expected: Command{
				Command: "my args",
				Args:    []string{"some"},
			},
		},
		{
			Command: "my %3%",
			Args:    []string{"some", "args"},
			Expected: Command{
				Command: "my ",
				Args:    []string{"some", "args"},
			},
		},
		{
			Command: "my %123%",
			Args:    []string{"some", "args"},
			Expected: Command{
				Command: "my ",
				Args:    []string{"some", "args"},
			},
		},
		{
			Command: "my %1%",
			Args:    []string{},
			Expected: Command{
				Command: "my ",
				Args:    []string{},
			},
		},
		{
			Command: "my %-1%",
			Args:    []string{"a", "b"},
			Expected: Command{
				Command: "my b",
				Args:    []string{"a"},
			},
		},
		{
			Command: "my %-1%",
			Args:    []string{},
			Expected: Command{
				Command: "my ",
				Args:    []string{},
			},
		},
		{
			Command: "k get secret -n %1% %2% -ojsonpath='{.data.%3%}'",
			Args:    []string{"ns", "name", "pass"},
			Expected: Command{
				Command: "k get secret -n ns name -ojsonpath='{.data.pass}'",
				Args:    []string{},
			},
		},
		{
			Command: "my %1%%2%%3%",
			Args:    []string{"hell", "o-", "world", "yay"},
			Expected: Command{
				Command: "my hello-world",
				Args:    []string{"yay"},
			},
		},
		{
			Command: "%1% %1% %2%",
			Args:    []string{"hey", "you"},
			Expected: Command{
				Command: "hey hey you",
				Args:    []string{},
			},
		},
	} {
		result := GenerateCommand(testCase.Command, testCase.Args)
		if result.Command != testCase.Expected.Command || slices.Compare(result.Args, testCase.Expected.Args) != 0 {
			t.Errorf("Expected %#v but got %#v", testCase.Expected, result)
		}
	}
}
