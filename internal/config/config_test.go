package config

import (
	"errors"
	"path"
	"reflect"
	"slices"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type ConfigFile struct {
	Path    string
	Content string
}

func resetConfig(fs afero.Fs, c Config, configFile ConfigFile) afero.Fs {
	fs.MkdirAll(path.Dir(configFile.Path), 0o755)
	f, err := fs.Create(configFile.Path)
	if err != nil {
		panic(err)
	}
	if _, err := f.WriteString(configFile.Content); err != nil {
		panic(err)
	}
	f.Close()
	c.v.SetFs(fs)
	c.v.SetConfigFile(configFile.Path)
	c.v.ReadInConfig()
	return fs
}

type SetAliasTestCase struct {
	Name           []string
	Command        []string
	ExistingConfig ConfigFile
	ExpectedError  error
	ExpectedConfig string
}

func TestSetAlias(t *testing.T) {
	for _, testCase := range []SetAliasTestCase{
		{
			Name:    []string{"my", "alias"},
			Command: []string{"command"},
			ExistingConfig: ConfigFile{
				Path: "/tmp/config.yaml",
			},
			ExpectedConfig: "my:\n    alias: command\n",
		},
		{
			Name:    []string{"my", "alias2"},
			Command: []string{"command", "args"},
			ExistingConfig: ConfigFile{
				Path: "/tmp/config.yaml",
			},
			ExpectedConfig: "my:\n    alias2: command args\n",
		},
	} {
		c := Config{v: viper.New()}
		fs := afero.NewMemMapFs()
		resetConfig(fs, c, testCase.ExistingConfig)
		if err := c.SetAlias(testCase.Name, testCase.Command); err != testCase.ExpectedError {
			t.Errorf("Expected %#v but got %#v", testCase.ExpectedError, err)
		}
		if contents, err := afero.ReadFile(fs, testCase.ExistingConfig.Path); err != nil {
			panic(err)
		} else {
			if result := string(contents); result != testCase.ExpectedConfig {
				t.Errorf("Config contents differ. Expected:\n%#v\nBut got:\n%#v", testCase.ExpectedConfig, result)
			}
		}
	}

}

type LoadFilesTestCase struct {
	Configs         []ConfigFile
	ExpectedContent map[string]any
}

func TestLoadFiles(t *testing.T) {
	for _, testCase := range []LoadFilesTestCase{
		{
			Configs: []ConfigFile{
				{
					Path:    "/tmp/config1.yaml",
					Content: "my:\n    alias: command\n",
				},
			},
			ExpectedContent: map[string]any{
				"my": map[string]string{
					"alias": "command",
				},
			},
		},
		{
			Configs: []ConfigFile{
				{
					Path:    "/tmp/config1.yaml",
					Content: "my:\n    alias: command\n",
				},
				{
					Path:    "/home/user/.config/na/na.yaml",
					Content: "my\n    alias2: command2\n",
				},
				{
					Path:    "/root/another.yaml",
					Content: "another: command3",
				},
			},
			ExpectedContent: map[string]any{
				"my": map[string]string{
					"alias":  "command",
					"alias2": "command2",
				},
				"another": "command3",
			},
		},
	} {
		c := Config{v: viper.New()}
		fs := afero.NewMemMapFs()
		for _, configFile := range testCase.Configs {
			resetConfig(fs, c, configFile)
		}
		paths := make([]string, 0)
		for _, configFile := range testCase.Configs {
			paths = append(paths, configFile.Path)
		}
		c.LoadFiles(paths...)

		if settings := c.v.AllSettings(); reflect.DeepEqual(settings, testCase.ExpectedContent) {
			t.Errorf("Expected %#v but got %#v", testCase.ExpectedContent, settings)
		}
	}
}

type ListAliasesTestCase struct {
	Config   ConfigFile
	Prefix   []string
	Expected []Alias
}

func TestListAliases(t *testing.T) {
	for _, testCase := range []ListAliasesTestCase{
		{
			Config: ConfigFile{
				Path:    "/tmp/config.yaml",
				Content: "my:\n    alias: command\nanother: command3",
			},
			Prefix: []string{"my"},
			Expected: []Alias{
				{
					Name:    "my alias",
					Command: "command",
				},
			},
		},
		{
			Config: ConfigFile{
				Path:    "/tmp/config.yaml",
				Content: "my:\n    alias: command\n    alias2: command2",
			},
			Prefix: []string{"my", "alias2"},
			Expected: []Alias{
				{
					Name:    "my alias2",
					Command: "command2",
				},
			},
		},
		{
			Config: ConfigFile{
				Path:    "/tmp/config.yaml",
				Content: "my:\n    alias: command\n    alias2: command2",
			},
			Prefix: []string{"my"},
			Expected: []Alias{
				{
					Name:    "my alias",
					Command: "command",
				},
				{
					Name:    "my alias2",
					Command: "command2",
				},
			},
		},
	} {
		c := Config{v: viper.New()}
		fs := afero.NewMemMapFs()
		resetConfig(fs, c, testCase.Config)
		if aliases := c.ListAliases(testCase.Prefix...); !slices.Equal(aliases, testCase.Expected) {
			t.Errorf("Expected %#v but got %#v", testCase.Expected, aliases)
		}
	}
}

type UnsetAliasTestCase struct {
	Config         ConfigFile
	Key            []string
	ExpectedConfig string
	ExpectedError  error
}

func TestUnsetAlias(t *testing.T) {
	for _, testCase := range []UnsetAliasTestCase{
		{
			Config: ConfigFile{
				Path:    "/tmp/config.yaml",
				Content: "my:\n    alias: cmd\n    alias2: cmd2\n",
			},
			Key:            []string{"my", "alias"},
			ExpectedConfig: "my:\n    alias2: cmd2\n",
		},
		{
			Config: ConfigFile{
				Path:    "/tmp/config.yaml",
				Content: "my:\n    alias: cmd\n    alias2: cmd2\n",
			},
			Key:            []string{"my"},
			ExpectedConfig: "{}\n",
		},
		{
			Config: ConfigFile{
				Path:    "/tmp/config.yaml",
				Content: "my:\n    alias: cmd\n    alias2: cmd2\n",
			},
			Key:            []string{"another"},
			ExpectedConfig: "my:\n    alias: cmd\n    alias2: cmd2\n",
			ExpectedError:  ErrAliasNotFound,
		},
		{
			Config: ConfigFile{
				Path:    "/tmp/config.yaml",
				Content: "my:\n    alias: cmd1\n    alias2: cmd2\n",
			},
			Key:            []string{"my", "alias", "invalid"},
			ExpectedConfig: "my:\n    alias: cmd1\n    alias2: cmd2\n",
			ExpectedError:  ErrInvalidAliasKey{"my alias"},
		},
	} {
		c := Config{v: viper.New()}
		fs := afero.NewMemMapFs()
		resetConfig(fs, c, testCase.Config)
		if err := c.UnsetAlias(testCase.Key...); !errors.Is(err, testCase.ExpectedError) {
			t.Errorf("Expected %#v but got %#v", testCase.ExpectedError, err)
		}
		if contents, err := afero.ReadFile(fs, testCase.Config.Path); err != nil {
			panic(err)
		} else {
			if result := string(contents); result != testCase.ExpectedConfig {
				t.Errorf("Config contents differ. Expected:\n%#v\nBut got:\n%#v", testCase.ExpectedConfig, result)
			}
		}
	}
}

type UnsetAliasByPrefixTestCase struct {
	Config         ConfigFile
	Key            []string
	ExpectedConfig string
	ExpectedError  error
}

func TestUnsetAliasByPrefix(t *testing.T) {
	for _, testCase := range []UnsetAliasTestCase{
		{
			Config: ConfigFile{
				Path:    "/tmp/config.yaml",
				Content: "my:\n    alias: cmd\n    alias2: cmd2\n",
			},
			Key:            []string{"my", "alias"},
			ExpectedConfig: "{}\n",
		},
		{
			Config: ConfigFile{
				Path:    "/tmp/config.yaml",
				Content: "my:\n    alias: cmd\n    alias2: cmd2\n",
			},
			Key:            []string{"my"},
			ExpectedConfig: "{}\n",
		},
		{
			Config: ConfigFile{
				Path:    "/tmp/config.yaml",
				Content: "my:\n    alias: cmd\n    alias2: cmd2\n",
			},
			Key:            []string{"another"},
			ExpectedConfig: "my:\n    alias: cmd\n    alias2: cmd2\n",
			ExpectedError:  ErrAliasNotFound,
		},
	} {
		c := Config{v: viper.New()}
		fs := afero.NewMemMapFs()
		resetConfig(fs, c, testCase.Config)
		if err := c.UnsetAliasByPrefix(testCase.Key...); !errors.Is(err, testCase.ExpectedError) {
			t.Errorf("Expected %#v but got %#v", testCase.ExpectedError, err)
		}
		if contents, err := afero.ReadFile(fs, testCase.Config.Path); err != nil {
			panic(err)
		} else {
			if result := string(contents); result != testCase.ExpectedConfig {
				t.Errorf("Config contents differ. Expected:\n%#v\nBut got:\n%#v", testCase.ExpectedConfig, result)
			}
		}
	}
}

type GetAliasTestCase struct {
	Config        ConfigFile
	Key           []string
	ExpectedAlias Alias
	ExpectedError error
}

func TestGetAlias(t *testing.T) {
	for _, testCase := range []GetAliasTestCase{
		{
			Config: ConfigFile{
				Path:    "/tmp/config.yaml",
				Content: "my:\n    alias: command\n",
			},
			Key:           []string{},
			ExpectedAlias: Alias{},
			ExpectedError: ErrAliasNotFound,
		},
		{
			Config: ConfigFile{
				Path:    "/tmp/config.yaml",
				Content: "my:\n    alias: command\n",
			},
			Key:           []string{"my"},
			ExpectedAlias: Alias{},
			ExpectedError: ErrAliasNotFound,
		},
		{
			Config: ConfigFile{
				Path:    "/tmp/config.yaml",
				Content: "my:\n    alias: command\n",
			},
			Key:           []string{"my", "alias"},
			ExpectedAlias: Alias{Name: "my alias", Command: "command"},
			ExpectedError: nil,
		},
	} {
		c := Config{v: viper.New()}
		fs := afero.NewMemMapFs()
		resetConfig(fs, c, testCase.Config)
		if a, err := c.GetAlias(testCase.Key...); a != testCase.ExpectedAlias || err != testCase.ExpectedError {
			t.Errorf("Expected (%#v, %#v) but got (%#v, %#v)", testCase.ExpectedAlias, testCase.ExpectedError, a, err)
		}
	}
}
