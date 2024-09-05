package dbmanager

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Database DatabaseConfig `json:"database"`
}

type DatabaseConfig struct {
	Port     int    `json:"port"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

func InitDatabase() error {

	_, err := os.Stat("config.json")
	if err == nil {
		return nil
	}

	config := Config{
		Database: DatabaseConfig{
			Port:     8080,
			Host:     "localhost",
			User:     "root",
			Password: "password",
			DBName:   "mydb",
		},
	}

	data, err := json.MarshalIndent(&config, "", "  ")
	if err != nil {
		fmt.Println("Could not indent json")
		return err
	}

	err = os.WriteFile("config.json", data, 0644)
	if err != nil {
		fmt.Println("Could not write config.json")
		return err
	}
	return nil
}
