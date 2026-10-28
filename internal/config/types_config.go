package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Db_url           string `json:"db_url"`
	Current_Username string `json:"current_user_name"`
}

func Read() (*Config, error) {
	configFileLocation, err := getConfigPath()
	if err != nil {
		return &Config{}, err
	}

	file, err := os.Open(configFileLocation)
	if err != nil {
		return &Config{}, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	theConfig := Config{}
	err = decoder.Decode(&theConfig)
	if err != nil {
		return &Config{}, err
	}

	return &theConfig, nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.Current_Username = userName
	err := writeConfig(*cfg) // takes the dereferenced data, cfg is not mutated
	if err != nil {
		return err
	}

	return nil
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configFileLocation := filepath.Join(homeDir, configFileName)

	return configFileLocation, nil
}

func writeConfig(cfg Config) error {
	configFileLocation, err := getConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(configFileLocation) // creates the file if it doesn't exist
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}
