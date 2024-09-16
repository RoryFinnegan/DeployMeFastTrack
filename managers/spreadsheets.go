// Excelize Library Licensing Clause

/*
BSD 3-Clause License

Copyright (c) 2016-2024 The excelize Authors.
Copyright (c) 2011-2017 Geoffrey J. Teale
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of the copyright holder nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package manager

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

var Spreadsheet SpreadsheetConfig

var ssMutex sync.Mutex

func spreadsheetGetConfig() error {
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

	Spreadsheet = config.Spreadsheet

	return nil
}

// Create and return a spreadsheet
func InitSpreadsheet() *excelize.File {
	var spreadsheet *excelize.File

	err := spreadsheetGetConfig()
	if err != nil {
		fmt.Println(err)
	}

	//Check if a spreadsheet already exists
	spreadsheet, err = excelize.OpenFile(Spreadsheet.Path)
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
func UpdateSpreadsheet(spreadsheet *excelize.File, asset string, serial string, user string) {
	ssMutex.Lock()
	defer ssMutex.Unlock()

	count := findStart(spreadsheet)

	spreadsheet.SetCellValue("Sheet1", fmt.Sprintf("A%d", count), asset)
	spreadsheet.SetCellValue("Sheet1", fmt.Sprintf("B%d", count), serial)
	spreadsheet.SetCellValue("Sheet1", fmt.Sprintf("C%d", count), user)
	//Take the current time and add to the fourth column
	spreadsheet.SetCellValue("Sheet1", fmt.Sprintf("D%d", count), time.Now())

	if err := spreadsheet.SaveAs(Spreadsheet.Path); err != nil {
		fmt.Println(err)
	}

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
