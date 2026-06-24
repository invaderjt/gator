package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func ReadConfig() (Config, error) {
	file, err := getConfigFilePath()
	if err != nil {
		fmt.Println("Could not find config file")
	}

	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Could not read file")
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Println("Could not unmarshal json")
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	err := write(*c)
	if err != nil {
		fmt.Println("Error writing username to file")
		return err
	}

	fmt.Printf("Username: %s written to file\n", username)
	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Could not identify home directory")
		return "", err
	}

	file := home + "/.gatorconfig.json"
	return file, nil
}

func write(cfg Config) error {
	file, err := getConfigFilePath()
	if err != nil {
		fmt.Println("Could not find config file")
		return err
	}

	jsonData, err := json.Marshal(cfg)
	if err != nil {
		fmt.Println("Could not marshal json data")
		return err
	}

	err = os.WriteFile(file, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file")
		return err
	}
	fmt.Println("Write successful")
	return nil
}
