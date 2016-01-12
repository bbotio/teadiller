package excel

import (
	"teadiller/models"
	"errors"
	"teadiller/Godeps/_workspace/src/github.com/tealeg/xlsx"
	"strings"
	"strconv"
)

type ExcelItemDao struct {
	FilePath string
	items    []*models.Item
}

func (excel ExcelItemDao) GetById(id string) (*models.Item, error) {
	for _, item := range excel.items {
		if item.Id == id {
			return item, nil
		}
	}

	return nil, errors.New("Item with id = " + id + " not found.")
}

func (excel ExcelItemDao) GetAll() []*models.Item {
	if excel.items == nil {
		excel.ParseItems()
	}

	return excel.items
}

func (excel *ExcelItemDao) ParseItems() {
	xlFile, _ := xlsx.OpenFile(excel.FilePath)
	sheet, _ := getSheetByName(xlFile, "items")
	excel.items = parseItemsSheet(sheet)
}

func parseItemsSheet(sheet *xlsx.Sheet) []*models.Item {
	var items []*models.Item

	for rowNum, row := range sheet.Rows {
		if rowNum == 0 {
			continue
		}

		item := new(models.Item)
		for i, cell := range row.Cells {
			cell.Value = strings.TrimSpace(cell.Value)
			switch i {
			case 0:
				item.Id = cell.Value
			case 1:
				item.Name = cell.Value
			case 2:
				item.Desc = cell.Value
			case 3:
				item.PhotoPath = cell.Value
			case 4:
				item.Tags = strings.Fields(cell.Value)
			case 5:
				// TODO: how to store this filed in excel?
				// We can't define all available values for Volume/Count types
				item.Type = models.VolumeItem
			case 6:
				item.AdditionalFields = make(map[string]string)
				properties := strings.Fields(cell.Value)
				for _, property := range properties {
					if len(property) == 0 {
						continue
					}

					keyvalue := strings.Split(property, "=")
					if len(keyvalue) == 2 {
						item.AdditionalFields[keyvalue[0]] = keyvalue[1]
					}
				}
			case 7:
				item.Count, _ = strconv.ParseFloat(cell.Value, 64)
			}
		}

		if (item.Id != "") {
			items = append(items, item)
		}
	}

	return items
}

func getSheetByName(file *xlsx.File, sheetName string) (*xlsx.Sheet, error) {
	for _, sheet := range file.Sheets {
		if sheet.Name == sheetName {
			return sheet, nil
		}
	}

	return nil, errors.New("Sheet with name = " + sheetName + " not found.")
}
