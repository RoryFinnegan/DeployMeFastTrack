package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

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

		// Check if the user pressed the Tab key
		if buffer[0] == '\t' {
			return strings.TrimSpace(string(input))
		}

		// Append the character to the input slice
		input = append(input, buffer[0])

		// Check if the user pressed Enter (newline)
		if buffer[0] == '\n' {
			return strings.TrimSpace(string(input))
		}
	}
	return strings.TrimSpace(string(input))
}

func main() {

	var f *excelize.File
	f, err := excelize.OpenFile("deploymebook.xlsx")
	if err != nil {
		fmt.Println("Book not found... making one now!")
		f = excelize.NewFile()

		err := f.SetColWidth("Sheet1", "A", "D", 50)
		if err != nil {
			fmt.Println("Could not make columns wide :(")
		}

		f.SetCellValue("Sheet1", "A1", "ASSET TAG")
		f.SetCellValue("Sheet1", "B1", "SERIAL")
		f.SetCellValue("Sheet1", "C1", "TECHNICIAN")
		f.SetCellValue("Sheet1", "D1", "TIME")
	}

	count := findStart(f)
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

		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", count), asset)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", count), serial)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", count), user)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", count), time.Now())

		if err := f.SaveAs("deploymebook.xlsx"); err != nil {
			fmt.Println(err)
		}
		count++
	}

}
