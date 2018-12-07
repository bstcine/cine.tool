package utils

import (
	"github.com/Luxurioust/excelize"
	"fmt"
	"github.com/tealeg/xlsx"
)

func ReadXLSX(path string) [][]string{
	xlsxFile, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
	}
	name := xlsxFile.GetSheetName(1)
	rows := xlsxFile.GetRows(name);
	if rows == nil || len(rows) == 0 {
		return nil
	}
	newRows := rows[1:]
	return newRows
}

func CreateXLMS(path string, keys []string) (*excelize.File, error) {
	xlsx := excelize.NewFile()
	for index,key := range keys {
		axis := getAxis(index)
		xlsx.SetCellStr("Sheet1", axis+"1", key)
	}
	err := xlsx.SaveAs(path)
	return xlsx, err
}

func AddElement(path string, values []string) error {
	file,err := xlsx.OpenFile(path)
	sheet := file.Sheets[0]
	row := sheet.AddRow()
	for _,value := range values{
		cell := row.AddCell()
		cell.Value = value
	}
	err = file.Save(path)
	return err
}

func getAxis(index int) string{
	b := byte(index+65)
	return  string(b)
}