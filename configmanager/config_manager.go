package configmanager

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type ConfigManager struct {
	koanf  *koanf.Koanf
	option Option
}

type Option struct {
	Delimeter string
	Path      string
}

func New(opt Option) *ConfigManager {
	if opt.Delimeter == "" {
		opt.Delimeter = "."
	}

	return &ConfigManager{
		koanf:  koanf.New(opt.Delimeter),
		option: opt,
	}
}

func (c *ConfigManager) throwErrorIfNotExists(path string) error {
	if !c.Exists(path) {
		return fmt.Errorf("config path: %s not found", path)
	}

	return nil
}

func (c *ConfigManager) Exists(path string) bool {
	return c.koanf.Exists(path)
}

func (c *ConfigManager) GetString(path string) (string, error) {
	err := c.throwErrorIfNotExists(path)
	if err != nil {
		return "", err
	}

	return c.koanf.String(path), nil
}

func (c *ConfigManager) GetBool(path string) (bool, error) {
	err := c.throwErrorIfNotExists(path)
	if err != nil {
		return false, err
	}

	return c.koanf.Bool(path), nil
}

func (c *ConfigManager) GetInt(path string) (int, error) {
	err := c.throwErrorIfNotExists(path)
	if err != nil {
		return 0, err
	}

	return c.koanf.Int(path), nil
}

func (c *ConfigManager) GetFloat64(path string) (float64, error) {
	err := c.throwErrorIfNotExists(path)
	if err != nil {
		return 0, err
	}

	return c.koanf.Float64(path), nil
}

func (c *ConfigManager) GetMapString(path string) (map[string]string, error) {
	err := c.throwErrorIfNotExists(path)
	if err != nil {
		return map[string]string{}, err
	}

	return c.koanf.StringMap(path), nil
}

func (c *ConfigManager) GetAll() (map[string]interface{}, error) {
	return c.koanf.All(), nil
}

func (c *ConfigManager) Load(configMap interface{}) error {
	err := c.loadFromFile()
	if err != nil {
		return err
	}

	err = c.loadFromEnvVar()
	if err != nil {
		return err
	}

	err = c.unmarshal(configMap)
	if err != nil {
		return err
	}

	return nil
}

func (c *ConfigManager) loadFromFile() error {
	if c.option.Path == "" {
		return nil
	}

	subPath := strings.Split(c.option.Path, ".")
	ext := subPath[len(subPath)-1]
	var parser koanf.Parser
	switch ext {
	case "json":
		parser = json.Parser()
	case "yaml", "yml":
		parser = yaml.Parser()
	case "env":
		parser = dotenv.Parser()
	case "toml":
		parser = toml.Parser()
	default:
		return fmt.Errorf("unsupported file type: %s", ext)
	}

	err := c.koanf.Load(file.Provider(c.option.Path), parser)
	if err != nil {
		return err
	}

	return nil
}

func (c *ConfigManager) loadFromEnvVar() error {
	return c.koanf.Load(env.Provider("", "_", func(s string) string {
		return strings.ReplaceAll(s, "_", c.option.Delimeter)
	}), nil)
}

func (c *ConfigManager) unmarshal(configMap interface{}) error {
	return c.koanf.UnmarshalWithConf("", configMap, koanf.UnmarshalConf{Tag: "conf", FlatPaths: true})
}
