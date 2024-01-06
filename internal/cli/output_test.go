package cli

import (
	"testing"

	"github.com/moyiz/na/internal/config"
)

func TestSimpleOutput(t *testing.T) {
	for _, testCase := range []struct {
		Aliases     []config.Alias
		ExpectedOut string
		ExpectedErr error
	}{
		{
			Aliases: []config.Alias{
				{
					Name:    "alias1",
					Command: "cmd1",
				},
				{
					Name:    "alias2",
					Command: "cmd2",
				},
				{
					Name:    "my alias1",
					Command: "my cmd1",
				},
				{
					Name:    "my alias2",
					Command: "my cmd2",
				},
			},
			ExpectedOut: "alias1 -- cmd1\nalias2 -- cmd2\nmy alias1 -- my cmd1\nmy alias2 -- my cmd2\n",
		},
		{
			Aliases:     []config.Alias{},
			ExpectedOut: "",
		},
	} {
		if out, err := simpleOutput(testCase.Aliases); out != testCase.ExpectedOut || err != testCase.ExpectedErr {
			t.Errorf("Simple output differ. Expected (%#v, %#v) but got (%#v, %#v)", testCase.ExpectedOut, testCase.ExpectedErr, out, err)
		}
	}
}

func TestCmdOutput(t *testing.T) {
	for _, testCase := range []struct {
		Aliases     []config.Alias
		ExpectedOut string
		ExpectedErr error
	}{
		{
			Aliases: []config.Alias{
				{
					Name:    "alias1",
					Command: "cmd1",
				},
				{
					Name:    "alias2",
					Command: "cmd2",
				},
				{
					Name:    "my alias1",
					Command: "my cmd1",
				},
				{
					Name:    "my alias2",
					Command: "my cmd2",
				},
			},
			ExpectedOut: "na add alias1 -- cmd1\nna add alias2 -- cmd2\nna add my alias1 -- my cmd1\nna add my alias2 -- my cmd2\n",
		},
		{
			Aliases:     []config.Alias{},
			ExpectedOut: "",
		},
	} {
		if out, err := cmdOutput(testCase.Aliases); out != testCase.ExpectedOut || err != testCase.ExpectedErr {
			t.Errorf("Command output differ. Expected (%#v, %#v) but got (%#v, %#v)", testCase.ExpectedOut, testCase.ExpectedErr, out, err)
		}
	}
}

func TestYamlOutput(t *testing.T) {
	for _, testCase := range []struct {
		Aliases     []config.Alias
		ExpectedOut string
		ExpectedErr error
	}{
		{
			Aliases:     []config.Alias{},
			ExpectedOut: "",
		},
		{
			Aliases: []config.Alias{
				{
					Name:    "alias1",
					Command: "cmd1",
				},
			},
			ExpectedOut: "---\nalias1: cmd1\n",
		},
		{
			Aliases: []config.Alias{
				{
					Name:    "my alias1",
					Command: "my cmd1",
				},
				{
					Name:    "my alias2",
					Command: "my cmd2",
				},
			},
			ExpectedOut: "---\nmy:\n    alias1: my cmd1\n    alias2: my cmd2\n",
		},
	} {
		if out, err := yamlOutput(testCase.Aliases); out != testCase.ExpectedOut || err != testCase.ExpectedErr {
			t.Errorf("Command output differ. Expected (%#v, %#v) but got (%#v, %#v)", testCase.ExpectedOut, testCase.ExpectedErr, out, err)
		}
	}
}
