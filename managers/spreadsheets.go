package manager

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

func InitSpreadsheet() *excelize.File {
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

func UpdateSpreadsheet(spreadsheet *excelize.File, count int, asset string, serial string, user string) {
	spreadsheet.SetCellValue("Sheet1", fmt.Sprintf("A%d", count), asset)
	spreadsheet.SetCellValue("Sheet1", fmt.Sprintf("B%d", count), serial)
	spreadsheet.SetCellValue("Sheet1", fmt.Sprintf("C%d", count), user)
	spreadsheet.SetCellValue("Sheet1", fmt.Sprintf("D%d", count), time.Now())

	if err := spreadsheet.SaveAs("deploymebook.xlsx"); err != nil {
		fmt.Println(err)
	}

}
