package common

import (
	"errors"

	"github.com/tealeg/xlsx"
)

func ReadXlsx(file_name string) (error, [][]string) {
	var xls_list [][]string
	xl_file, err := xlsx.OpenFile(file_name)
	if err != nil {
		return err, xls_list
	}
	for _, sheet := range xl_file.Sheets {
		for _, row := range sheet.Rows {
			tmp_file := []string{}

			for _, cell := range row.Cells {
				text := cell.String()
				tmp_file = append(tmp_file[0:], text)
			}
			xls_list = append(xls_list[0:], tmp_file)
		}
	}
	return err, xls_list
}

func WriteXlsx(file_name string, content [][]string) error {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return err
	}
	if len(content) < 1 {
		err = errors.New("content null")
		return err
	}
	for _, val_arr := range content {
		if len(val_arr) < 1 {
			continue
		}
		var row *xlsx.Row
		row = sheet.AddRow()
		for _, val := range val_arr {
			cell = row.AddCell()
			cell.Value = val
		}
	}
	err = file.Save(file_name)
	return err
}

func AppendWriteXlsx(file_name string, content [][]string) error {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var cell *xlsx.Cell
	var err error
	file, err = xlsx.OpenFile(file_name)
	if err != nil {
		panic(err)
	}
	if len(content) < 1 {
		err = errors.New("content null")
		return err
	}
	sheet = file.Sheets[0]
	for _, val_arr := range content {
		if len(val_arr) < 1 {
			continue
		}
		var row *xlsx.Row
		row = sheet.AddRow()
		for _, val := range val_arr {
			cell = row.AddCell()
			cell.Value = val
		}
	}
	err = file.Save(file_name)
	return err
}
