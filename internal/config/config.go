package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() string {
	configFilePath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return configFilePath + "/" + configFileName
}

func Read() (*Config, error) {
	configJSON, err := os.ReadFile(getConfigFilePath())
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func Write(config Config) error {
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(getConfigFilePath(), configJSON, 0644)
}

func (config *Config) SetUser(username string) error {
	config.CurrentUserName = username
	err := Write(*config)
	if err != nil {
		return err
	}
	return nil
}
