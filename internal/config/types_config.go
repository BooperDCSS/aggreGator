package config

import (
	"encoding/json"
	"fmt"
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

	data, err := os.ReadFile(configFileLocation)
	if err != nil {
		return &Config{}, err
	}

	var theConfig Config

	if err := json.Unmarshal(data, &theConfig); err != nil {
		return &Config{}, err
	}

	return &theConfig, nil
}

func (cfg *Config) SetUser(userName string) {
	cfg.Current_Username = userName
	err := writeConfig(*cfg)
	if err != nil {
		fmt.Printf("Error writing file: %v", err)
	}

}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error determining user home directory: %w", err)
	}

	configFileLocation := filepath.Join(homeDir, configFileName)

	return configFileLocation, nil
}

func writeConfig(cfg Config) error {
	configFileLocation, err := getConfigPath()
	if err != nil {
		return err
	}

	newConfigData, err := json.Marshal(cfg)
	if err != nil {
		fmt.Printf("Error marshaling data: %v", err)
	}
	if err := os.WriteFile(configFileLocation, newConfigData, 0666); err != nil {
		return fmt.Errorf("Error writing config with new userName: %w", err)
	}
	return nil
}
