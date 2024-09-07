package main

import (
	manager "DeployMeFastTrack/managers"

	"fmt"
	"os"
	"strings"
)

func readScan() string {
	//Create a buffer and a array to hold all values
	var input []byte
	buffer := make([]byte, 1)

	//Read until end of file is reached
	for {
		_, err := os.Stdin.Read(buffer)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		//Leave after a tab
		if buffer[0] == '\t' {
			return strings.TrimSpace(string(input))
		}

		input = append(input, buffer[0])

		//Leave after a newline
		if buffer[0] == '\n' {
			return strings.TrimSpace(string(input))
		}
	}
	return strings.TrimSpace(string(input))

}

// Program Start
func main() {
	//Initialize all the setup data
	spreadsheet := manager.InitSpreadsheet()
	manager.InitUserlist()
	manager.InitDatabase()
	manager.DBGetConfig()

	//Find the start of the spreadsheet to start inserting new values
	count := manager.FindStart(spreadsheet)
	var asset, serial, user string
	for {
		//Reset the asset, serial and user
		asset = ""
		serial = ""
		user = ""

		fmt.Println("Enter Asset Tag: ")
		for asset == "" {
			asset = readScan()
		}

		fmt.Println("Enter Serial Number: ")
		for serial == "" {
			serial = readScan()
		}

		fmt.Println("Enter Name: ")
		for user == "" {
			user = readScan()
		}
		user = manager.GetNameFromId(user)

		if manager.DBConfig.Driver != "None" {
			manager.InsertDatabaseRow(asset, serial, user)
		}
		manager.UpdateSpreadsheet(spreadsheet, count, asset, serial, user)
		count++
	}

}
