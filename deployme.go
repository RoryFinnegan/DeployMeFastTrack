package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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
	var input []byte
	buffer := make([]byte, 1)

	for {
		_, err := os.Stdin.Read(buffer)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		if buffer[0] == '\t' {
			return strings.TrimSpace(string(input))
		}

		input = append(input, buffer[0])

		if buffer[0] == '\n' {
			return strings.TrimSpace(string(input))
		}
	}
	return strings.TrimSpace(string(input))

	// prev_state, err := term.MakeRaw(int(os.Stdin.Fd()))
	// if err != nil {
	// 	fmt.Println("Terminal could not enter raw mode. Error:\n", err)
	// }

	// defer term.Restore(int(os.Stdin.Fd()), prev_state)

	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// go func() {
	// 	sig := <-sigChan
	// 	fmt.Printf("Received signal: %v. Restoring terminal and exiting.\n", sig)
	// 	term.Restore(int(os.Stdin.Fd()), prev_state)
	// 	os.Exit(0)
	// }()

	// for {
	// 	_, err = os.Stdin.Read(buffer)
	// 	if err != nil {
	// 		fmt.Println("Could not read buffer: \n", err)
	// 	}

	// 	fmt.Printf("%q", buffer[0])

	// 	if buffer[0] == '\n' || buffer[0] == '\t' || buffer[0] == '\r' {
	// 		os.Stdout.Write([]byte("\n"))
	// 		return string(input)
	// 	}

	// 	if buffer[0] == '\x03' {
	// 		term.Restore(int(os.Stdin.Fd()), prev_state)
	// 		os.Exit(0)
	// 	}

	// 	input = append(input, buffer[0])
	// }

}

func init_spreadsheet() *excelize.File {
	var spreadsheet *excelize.File
	spreadsheet, err := excelize.OpenFile("deploymebook.xlsx")
	if err != nil {
		fmt.Println("Book not found... making one now!")
		spreadsheet = excelize.NewFile()

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

	if _, err := os.Stat("userlist.json"); err == nil {
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

	encoder := json.NewEncoder(file)
	err = encoder.Encode(users)
	if err != nil {
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

	count := findStart(spreadsheet)
	// scanner := bufio.NewScanner(os.Stdin)
	var asset, serial, user string
	for {
		asset = ""
		serial = ""
		user = ""

		fmt.Println("Enter Asset Tag: ")
		// if scanner.Scan() {
		// 	asset = scanner.Text()
		// }
		// if err := scanner.Err(); err != nil {
		// 	fmt.Fprintln(os.Stderr, "Error reading input:", err)
		// }
		for asset == "" {
			asset = readScan()
		}

		fmt.Println("Enter Serial Number: ")
		// if scanner.Scan() {
		// 	serial = scanner.Text()
		// }
		// if err := scanner.Err(); err != nil {
		// 	fmt.Fprintln(os.Stderr, "Error reading input:", err)
		// }
		for serial == "" {
			serial = readScan()
		}

		fmt.Println("Enter Name: ")
		// if scanner.Scan() {
		// 	user = scanner.Text()
		// }
		// if err := scanner.Err(); err != nil {
		// 	fmt.Fprintln(os.Stderr, "Error reading input:", err)
		// }

		for user == "" {
			user = readScan()
		}
		user = get_name_from_id(user)

		spreadsheet.SetCellValue("Sheet1", fmt.Sprintf("A%d", count), asset)
		spreadsheet.SetCellValue("Sheet1", fmt.Sprintf("B%d", count), serial)
		spreadsheet.SetCellValue("Sheet1", fmt.Sprintf("C%d", count), user)
		spreadsheet.SetCellValue("Sheet1", fmt.Sprintf("D%d", count), time.Now())

		if err := spreadsheet.SaveAs("deploymebook.xlsx"); err != nil {
			fmt.Println(err)
		}
		count++
	}

}
