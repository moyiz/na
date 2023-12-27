package config

import (
	"bytes"
	"encoding/json"
	"path"
	"sort"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	v *viper.Viper
}

var c Config

func init() {
	c = Config{v: viper.New()}
}

type Alias struct {
	Name    string
	Command string
}

func LoadFiles(filePath ...string) map[string]any {
	return c.LoadFiles(filePath...)
}

func (c *Config) LoadFiles(filePath ...string) map[string]any {
	for i, p := range filePath {
		c.v.SetConfigFile(path.Clean(p))
		if i == 0 {
			c.v.ReadInConfig()
		} else {
			c.v.MergeInConfig()
		}
	}
	return c.v.AllSettings()
}

func SetAlias(key []string, command []string) error {
	return c.SetAlias(key, command)
}

func (c *Config) SetAlias(name []string, command []string) error {
	c.v.Set(strings.Join(name, "."), strings.Join(command, " "))
	return c.v.WriteConfig()
}

func UnsetAlias(key ...string) error {
	return c.UnsetAlias(key...)
}

func (c *Config) UnsetAlias(key ...string) error {
	var parent map[string]interface{}
	var keyIsMap, keyExists bool
	keySize := len(key)
	settings := c.v.AllSettings()
	aliasWalker := settings
	for i, k := range key {
		parent = aliasWalker
		_, keyExists = aliasWalker[k]
		aliasWalker, keyIsMap = aliasWalker[k].(map[string]interface{})
		if !keyExists {
			return ErrAliasNotFound
		} else if !keyIsMap && i < keySize-1 {
			return ErrInvalidAliasKey{strings.Join(key[:i+1], " ")}
		} else if i == keySize-1 {
			break
		}
	}
	delete(parent, key[keySize-1])
	encoded, _ := json.MarshalIndent(settings, "", " ")
	c.v.ReadConfig(bytes.NewReader(encoded))
	return c.v.WriteConfig()
}

func UnsetAliasByPrefix(prefix ...string) error {
	return c.UnsetAliasByPrefix(prefix...)
}

func (c *Config) UnsetAliasByPrefix(prefix ...string) error {
	var found bool
	keyPrefix := strings.Join(prefix, ".")
	for _, key := range c.v.AllKeys() {
		if strings.HasPrefix(key, keyPrefix) {
			keyParts := strings.Split(key, ".")
			parts := len(keyParts)
			keyPath := strings.Join(keyParts[:parts-1], ".")
			delete(c.v.Get(keyPath).(map[string]interface{}), keyParts[parts-1])
			found = true
		}
	}
	if !found {
		return ErrAliasNotFound
	}
	return c.v.WriteConfig()
}

func Write() error {
	return c.Write()
}

func (c *Config) Write() error {
	return c.v.WriteConfig()
}

func ListAliases(key ...string) []Alias {
	return c.ListAliases(key...)
}

func (c *Config) ListAliases(key ...string) []Alias {
	keys := c.v.AllKeys()
	sort.Strings(keys)

	aliases := make([]Alias, 0)
	aliasPrefix := strings.Join(key, ".")
	for _, k := range keys {
		if strings.HasPrefix(k, aliasPrefix+".") || k == aliasPrefix {
			aliases = append(aliases, Alias{
				Name:    strings.ReplaceAll(k, ".", " "),
				Command: c.v.GetString(k),
			})
		}
	}

	return aliases
}

func ListAliasesByPrefix(prefix ...string) []Alias {
	return c.ListAliasesByPrefix(prefix...)
}

func (c *Config) ListAliasesByPrefix(prefix ...string) []Alias {
	keys := c.v.AllKeys()
	sort.Strings(keys)

	aliases := make([]Alias, 0)
	aliasPrefix := strings.Join(prefix, ".")
	for _, k := range keys {
		if strings.HasPrefix(k, aliasPrefix) {
			aliases = append(aliases, Alias{
				Name:    strings.ReplaceAll(k, ".", " "),
				Command: c.v.GetString(k),
			})
		}
	}

	return aliases
}

func GetAlias(part ...string) (Alias, error) {
	return c.GetAlias(part...)
}

func (c *Config) GetAlias(part ...string) (Alias, error) {
	key := strings.Join(part, ".")
	if command := c.v.GetString(key); command == "" {
		return Alias{}, ErrAliasNotFound
	} else {
		return Alias{Name: strings.ReplaceAll(key, ".", " "), Command: command}, nil
	}
}
