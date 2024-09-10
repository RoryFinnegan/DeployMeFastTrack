package manager

import (
	"encoding/json"
	"fmt"
	"os"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Struct used to hold the top layer of all data
// taken from config.json
type Config struct {
	Database DatabaseConfig `json:"database"`
	Mail     MailConfig     `json:"mail"`
}

// Holds the settings for connecting to a database
type DatabaseConfig struct {
	//The driver for connecting to a database i.e. "mysql", "postgresql", "sqlite"
	Driver string `json:"Driver"`
	//Port where the database is opened
	Port string `json:"port"`
	//The hostname for the database
	Host string `json:"host"`
	//The database user we are connecting as
	User string `json:"user"`
	//The database password for the given user
	Password string `json:"password"`
	//The name of the database we are connecting to
	DBName string `json:"dbname"`
}

type MailConfig struct {
	Sender   string `json:"Sender"`
	Password string `json:"Password"`
	Receiver string `json:"Receiver"`
	Server   string `json:"Server"`
	Port     string `json:"Port"`
}

// Create default database information for the config
func InitDatabase() error {
	//If the config exists already, leave
	_, err := os.Stat("config.json")
	if err == nil {
		return nil
	}

	//Prefill the Config with default information, does not enable database by defualt
	config := Config{
		Database: DatabaseConfig{
			Driver:   "None",
			Port:     "8080",
			Host:     "localhost",
			User:     "root",
			Password: "password",
			DBName:   "mydb",
		},
		Mail: MailConfig{
			Sender:   "sender@example.com",
			Password: "password",
			Receiver: "receiver@example.com",
			Server:   "mail.example.com",
			Port:     "25",
		},
	}

	//Give human-readable indentation to the database
	data, err := json.MarshalIndent(&config, "", "  ")
	if err != nil {
		fmt.Println("Could not indent json")
		return err
	}

	//Write the config
	err = os.WriteFile("config.json", data, 0644)
	if err != nil {
		fmt.Println("Could not write config.json")
		return err
	}

	return nil
}

// Initialize the userlist with a default user and some
// basic formatting
func InitUserlist() error {
	_, err := os.Stat("userlist.json")
	if err == nil {
		return nil
	}

	file, err := os.Create("userlist.json")
	if err != nil {
		return err
	}
	defer file.Close()

	users := []User{
		{ID: "1", Name: "TestAccount"},
	}

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		fmt.Println("Could not indent JSON")
		return err
	}

	err = os.WriteFile("userlist.json", data, 0644)
	if err != nil {
		fmt.Println("Could not write userlist.json")
		return err
	}

	return nil
}

// Take the map of users from userlist.json
func get_userlist() ([]User, error) {
	file, err := os.Open("userlist.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []User
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// search for the given id from the list of users
func GetNameFromId(id string) string {
	users, err := get_userlist()
	if err != nil {
		return string(id)
	}

	for _, user := range users {
		if user.ID == id {
			return user.Name
		}
	}

	return string(id)
}
