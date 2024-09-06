package manager

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Config struct {
	Database DatabaseConfig `json:"database"`
}

type DatabaseConfig struct {
	Driver   string `json:"Driver"`
	Port     int    `json:"port"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

var DBConfig DatabaseConfig
var db *sql.DB

func InitDatabase() error {

	_, err := os.Stat("config.json")
	if err == nil {
		return nil
	}

	config := Config{
		Database: DatabaseConfig{
			Driver:   "None",
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

func DBGetConfig() error {
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}

	DBConfig = config.Database

	return nil

}

func dbConnect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=prefer",
		DBConfig.Host, DBConfig.Port, DBConfig.User, DBConfig.Password, DBConfig.DBName)
	driver := DBConfig.Driver

	db, err := sql.Open(driver, psqlInfo)
	if err != nil {
		return nil, err
	}

	fmt.Println(driver)
	fmt.Println(psqlInfo)
	err = db.Ping()
	if err != nil {

		return nil, err
	}
	return db, nil

}

func InsertDatabaseRow(asset string, serial string, name string) error {
	db, err := dbConnect()
	if err != nil {
		fmt.Println(err)
		return err
	}

	stmt, err := db.Prepare("INSERT INTO assets (asset, serial, name) VALUES ($1, $2, $3)")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(asset, serial, name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
