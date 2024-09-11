package manager

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/go-sql-driver/mysql"
)

// Global variable to hold database configuration
var DBConfig DatabaseConfig

var dbMutex sync.Mutex

// Convert config information regarding the database into a struct
// then transfer information into global variable DBConfig
func DBGetConfig() error {
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	//Decode the json file into a GoLang struct
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}

	DBConfig = config.Database

	return nil

}

// Returns a connection to the database.
// NOTE: Any function calling this should make sure to
// close the database
func dbConnect() (*sql.DB, error) {
	//Provide information to the mysql config, currently
	//only passes information as needed for given environment
	cfg := mysql.Config{
		User:   DBConfig.User,
		Passwd: DBConfig.Password,
		Net:    "tcp",
		Addr:   DBConfig.Host + ":" + DBConfig.Port,
		DBName: DBConfig.DBName,
	}

	//Transform the config data into a suitable address for the database
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Println(err)
		return db, err
	}

	//Ping the database, makes sure we are connected
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return db, err
	}

	return db, nil
}

// Inserts a row into the database
// asset string: The asset tag of the device being added
// serial string: The serial number or product number of device being added
// name string: The hash or name of the employee registering the product
// All above are stored in the form of VARCHAR(255)
func InsertDatabaseRow(asset string, serial string, name string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	//Receive a connection to the DB
	db, err := dbConnect()
	if err != nil {
		return err
	}
	defer db.Close()

	//Create a prepared statement, use the current time on the SQL server
	stmt, err := db.Prepare("INSERT INTO assets (asset, snum, technician, date_added) VALUES (?, ?, ?, now())")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()

	//Execute the statement
	_, err = stmt.Exec(asset, serial, name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
