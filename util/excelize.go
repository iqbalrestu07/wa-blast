package util

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func OpenExcel(pathFile string) ([][]string, error) {
	f, err := excelize.OpenFile(pathFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")

	return rows, err
}

// GenerateExcel ...
func GenerateExcel(sheetName string, data map[string]string) *excelize.File {

	f := excelize.NewFile()

	f.NewSheet(sheetName)

	for k, v := range data {
		f.SetCellValue(sheetName, k, v)
	}

	return f
}
