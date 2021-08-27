package main

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type ConfigStruct struct {
	// Debug
	DebugMode bool `yaml:"DebugMode"`

	// Server
	VersionString string `yaml:"VersionString"`
	HostAndPort   string `yaml:"HostAndPort"`
	BaseURL       string `yaml:"BaseURL"`
	APIBase       string `yaml:"APIBase"`
	HTTPS         bool   `yaml:"HTTPS"`

	// Secret
	SecretFilename string

	// Database and Reddis
	DatabaseConnRetries int    `yaml:"DatabaseConnRetries"`
	DatabaseAddr        string `yaml:"DatabaseAddr"`
	DatabaseDatabase    string `yaml:"DatabaseDatabase"`
	DatabaseUser        string `yaml:"DatabaseUser"`
	DatabasePassword    string `yaml:"DatabasePassword"`
	RedisAddr           string `yaml:"RedisAddr"`

	// Users and Invites
	MaxUsernameSize int    `yaml:"MaxUsernameSize"`
	AdminUsername   string `yaml:"AdminUsername"`
	FirstUsername   string `yaml:"FirstUsername"`
	DefaultPassword string `yaml:"DefaultPassword"`
	InviteSize      int    `yaml:"InviteSize"`

	// Bots
	BotHearthBeatPeriod time.Duration `yaml:"BotHearthBeatPeriod"`
	BotHearthBeatExpire time.Duration `yaml:"BotHearthBeatExpire"`
	BotHearthBeatDead   time.Duration `yaml:"BotHearthBeatDead"`
	BotChannelLag       int           `yaml:"BotChannelLag"`

	// Templating
	TemplatesDir string `yaml:"TemplatesDir"`

	// Pages
	PageDefaultLimit  int `yaml:"PageDefaultLimit"`
	PageDefaultOffset int `yaml:"PageDefaultOffset"`
	PageDefaultSize   int `yaml:"PageDefaultSize"`
	PageDefaultNum    int `yaml:"PageDefaultNum"`
}

var Config *ConfigStruct

func init() {
	Config = new(ConfigStruct)

	if err := loadConfigFile(); err != nil {
		panic(err)
	}
}

func loadConfigFile() error {
	f, err := os.Open("src/config.yml")
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(Config)
	if err != nil {
		return err
	}

	return nil
}
