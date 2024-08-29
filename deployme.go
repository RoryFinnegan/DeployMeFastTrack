package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

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
	}

	scanner := bufio.NewScanner(os.Stdin)
	var asset, serial, user string
	for {

		fmt.Println("Enter Asset Tag: ")
		if scanner.Scan() {
			asset = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
		}

		fmt.Println("Enter Serial Number: ")
		if scanner.Scan() {
			serial = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
		}

		fmt.Println("Enter Name: ")
		if scanner.Scan() {
			user = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
		}

		f.SetCellValue("Sheet1", "A1", asset)
		f.SetCellValue("Sheet1", "B1", serial)
		f.SetCellValue("Sheet1", "C1", user)
		f.SetCellValue("Sheet1", "D1", time.Now())

		if err := f.SaveAs("deploymebook.xlsx"); err != nil {
			fmt.Println(err)
		}
	}

}
