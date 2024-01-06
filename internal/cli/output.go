package cli

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/moyiz/na/internal/config"
	"gopkg.in/yaml.v3"
)

type OutputFormat string

const (
	OutputFormatSimple = OutputFormat("simple")
	OutputFormatCmd    = OutputFormat("cmd")
	OutputFormatYaml   = OutputFormat("yaml")
)

func AllOutputFormats() []OutputFormat {
	return []OutputFormat{
		OutputFormatSimple,
		OutputFormatCmd,
		OutputFormatYaml,
	}
}

func (o *OutputFormat) String() string {
	return string(*o)
}

func (o *OutputFormat) Set(v string) error {
	switch {
	case slices.Contains(AllOutputFormats(), OutputFormat(v)):
		*o = OutputFormat(v)
		return nil
	default:
		return errors.New("invalid format: " + v)
	}
}

func (o *OutputFormat) Type() string {
	return "format"
}

var OutputFormatToFunc = map[OutputFormat](func([]config.Alias) (string, error)){
	OutputFormatSimple: simpleOutput,
	OutputFormatCmd:    cmdOutput,
	OutputFormatYaml:   yamlOutput,
}

func simpleOutput(aliases []config.Alias) (string, error) {
	builder := strings.Builder{}
	for _, alias := range aliases {
		builder.WriteString(fmt.Sprintln(alias.Name, "--", alias.Command))
	}
	return builder.String(), nil
}

func cmdOutput(aliases []config.Alias) (string, error) {
	builder := strings.Builder{}
	for _, alias := range aliases {
		builder.WriteString(fmt.Sprintln("na add", alias.Name, "--", alias.Command))
	}
	return builder.String(), nil
}

func yamlOutput(aliases []config.Alias) (string, error) {
	aliasesMap := make(map[string]interface{})
	for _, alias := range aliases {
		aliasWalker := aliasesMap
		aliasParts := strings.Fields(alias.Name)
		n := len(aliasParts)
		for i := 0; i < n-1; i++ {
			if _, ok := aliasWalker[aliasParts[i]]; !ok {
				aliasWalker[aliasParts[i]] = make(map[string]interface{})
			}
			aliasWalker = aliasWalker[aliasParts[i]].(map[string]interface{})
		}
		aliasWalker[aliasParts[n-1]] = alias.Command
	}
	builder := strings.Builder{}
	if yamlData, err := yaml.Marshal(aliasesMap); err == nil {
		if yamlStr := string(yamlData); yamlStr != "{}\n" {
			builder.WriteString(fmt.Sprintln("---"))
			builder.WriteString(yamlStr)
		}
	} else {
		return string(yamlData), err
	}
	return builder.String(), nil
}
