package manager

import (
	"fmt"
	"log"
	"time"

	"github.com/xuri/excelize/v2"
)

// Create and return a spreadsheet
func InitSpreadsheet() *excelize.File {
	var spreadsheet *excelize.File
	//Check if a spreadsheet already exists
	spreadsheet, err := excelize.OpenFile("deploymebook.xlsx")
	//If a spreadsheet doesn't exist, make one
	if err != nil {
		fmt.Println("Book not found... making one now!")
		spreadsheet = excelize.NewFile()

		err := spreadsheet.SetColWidth("Sheet1", "A", "D", 50)
		if err != nil {
			fmt.Println("Could not make columns wide :(")
		}
		//Create labeled columns
		spreadsheet.SetCellValue("Sheet1", "A1", "ASSET TAG")
		spreadsheet.SetCellValue("Sheet1", "B1", "SERIAL")
		spreadsheet.SetCellValue("Sheet1", "C1", "TECHNICIAN")
		spreadsheet.SetCellValue("Sheet1", "D1", "TIME")
	}
	return spreadsheet
}

// Add values to the spreadsheet
// spreadsheet: The file the spreadsheet is stored in
// count: The current row the add a value into
// asset: The asset tag to add
// serial: The serial/product number to add
// user: The technician to add
func UpdateSpreadsheet(spreadsheet *excelize.File, count int, asset string, serial string, user string) {
	spreadsheet.SetCellValue("Sheet1", fmt.Sprintf("A%d", count), asset)
	spreadsheet.SetCellValue("Sheet1", fmt.Sprintf("B%d", count), serial)
	spreadsheet.SetCellValue("Sheet1", fmt.Sprintf("C%d", count), user)
	//Take the current time and add to the fourth column
	spreadsheet.SetCellValue("Sheet1", fmt.Sprintf("D%d", count), time.Now())

	if err := spreadsheet.SaveAs("deploymebook.xlsx"); err != nil {
		fmt.Println(err)
	}

}

func FindStart(f *excelize.File) int {
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
