package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ProxyPort string
	AdminPort string
}

func LoadConfig(filename string) (Config, error) {
	// Load current configuration from JSON file
	file, err := os.Open(filename)
	if err != nil {
		file.Close()
		var config Config
		config.ProxyPort = ":8080"
		config.AdminPort = ":8008"
		error := SaveConfig(filename, config)
		if error != nil {
			// Log to file
		}
		return Config{}, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func SaveConfig(filename string, config Config) error {
	// Save current configuration to JSON file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}

	return nil
}
