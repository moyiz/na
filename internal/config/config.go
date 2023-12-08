package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"path"
	"sort"
	"strings"

	"github.com/spf13/viper"
)

func GetFromFiles(filePath ...string) *Config {
	for i, p := range filePath {
		viper.SetConfigFile(path.Clean(p))
		if i == 0 {
			viper.ReadInConfig()
		} else {
			viper.MergeInConfig()
		}
	}
	return &Config{aliases: viper.AllSettings()}
}

type Config struct {
	aliases map[string]interface{}
}

func (c *Config) Write() error {
	encoded, _ := json.MarshalIndent(c.aliases, "", " ")
	viper.ReadConfig(bytes.NewReader(encoded))
	return viper.WriteConfig()
}

func (c *Config) SetAlias(key []string, command []string) error {
	viper.Set(strings.Join(key, "."), strings.Join(command, " "))
	return viper.WriteConfig()
}

func (c *Config) UnsetAlias(key ...string) error {
	var parent map[string]interface{}
	var keyIsMap, keyExists bool
	keySize := len(key)
	aliasPointer := c.aliases
	for i, k := range key {
		parent = aliasPointer
		_, keyExists = aliasPointer[k]
		aliasPointer, keyIsMap = aliasPointer[k].(map[string]interface{})
		if !keyExists {
			return errors.New("not found")
		} else if !keyIsMap && i < keySize-1 {
			return errors.New("key is invalid. Did you mean `" + strings.Join(key[:i+1], " ") + "`?")
		} else if i == keySize-1 {
			break
		}
	}
	delete(parent, key[keySize-1])
	return c.Write()
}

type Alias struct {
	Name    string
	Command string
}

func (c *Config) ListAliases(prefix ...string) []Alias {
	return ListAliases(prefix...)
}

func ListAliases(prefix ...string) []Alias {
	keys := viper.AllKeys()
	sort.Strings(keys)

	aliases := make([]Alias, 0)
	aliasPrefix := strings.Join(prefix, ".")
	for _, k := range keys {
		if strings.HasPrefix(k, aliasPrefix) {
			aliases = append(aliases, Alias{
				Name:    strings.ReplaceAll(k, ".", " "),
				Command: viper.GetString(k),
			})
		}
	}

	return aliases
}

func (c *Config) GetAlias(part ...string) (Alias, error) {
	key := strings.Join(part, ".")
	if command := viper.GetString(key); command == "" {
		return Alias{}, errors.New("not found")
	} else {
		return Alias{Name: strings.ReplaceAll(key, ".", " "), Command: command}, nil
	}
}
