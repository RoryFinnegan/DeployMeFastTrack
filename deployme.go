package main

import (
	manager "DeployMeFastTrack/managers"

	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func findStart(f *excelize.File) int {
	for row := 1; ; row++ {
		// Get the value from column A for the current row
		cell, err := f.GetCellValue("Sheet1", fmt.Sprintf("A%d", row))
		if err != nil {
			log.Fatal(err)
		}

		// Check if the cell is empty
		if cell == "" {
			return row
		}

	}

}

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

func init_spreadsheet() *excelize.File {
	var spreadsheet *excelize.File
	spreadsheet, err := excelize.OpenFile("deploymebook.xlsx")
	if err != nil {
		//If a book does not exist create one
		fmt.Println("Book not found... making one now!")
		spreadsheet = excelize.NewFile()

		//TODO: Change program to allow customizable columns. Large task.
		err := spreadsheet.SetColWidth("Sheet1", "A", "D", 50)
		if err != nil {
			fmt.Println("Could not make columns wide :(")
		}

		spreadsheet.SetCellValue("Sheet1", "A1", "ASSET TAG")
		spreadsheet.SetCellValue("Sheet1", "B1", "SERIAL")
		spreadsheet.SetCellValue("Sheet1", "C1", "TECHNICIAN")
		spreadsheet.SetCellValue("Sheet1", "D1", "TIME")
	}
	return spreadsheet
}

func init_userlist() error {

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

func get_name_from_id(id string) string {
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

func main() {

	spreadsheet := init_spreadsheet()
	init_userlist()
	manager.InitDatabase()

	count := findStart(spreadsheet)
	var asset, serial, user string
	for {
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
		user = get_name_from_id(user)

		manager.UpdateSpreadsheet(spreadsheet, count, asset, serial, user)
		count++
	}

}
