package excel
import (
	"teadiller/Godeps/_workspace/src/github.com/tealeg/xlsx"
	"teadiller/models"
	"errors"
	"fmt"
	"time"
	"strings"
	"strconv"
)

type ExcelOrderDao struct {
	FilePath string
	orders   []*models.Order
}

func (excel ExcelOrderDao) GetById(id string) (*models.Order, error) {
	for _, order := range excel.orders {
		if order.Id == id {
			return order, nil
		}
	}
	return nil, errors.New("Order with id = " + id + " not found.")
}

func (excel ExcelOrderDao) GetAll() []*models.Order {
	if excel.orders == nil {
		excel.ParseOrders()
	}

	return excel.orders
}

func (excel ExcelOrderDao) Save(order *models.Order) {
	file, _ := xlsx.OpenFile(excel.FilePath)

	sheet, _ := getSheetByName(file, "orders")
	row := sheet.AddRow();

	row.AddCell().SetString(order.Id)
	row.AddCell().SetString(order.ItemId)
	row.AddCell().SetString(order.Count)
	row.AddCell().SetString(order.Buyer.Name)
	row.AddCell().SetString(order.Delivery.DeliveryType.String())
	row.AddCell().SetString(order.Delivery.Address)
	row.AddCell().SetString(fmt.Sprint(order.Datetime.Format(time.ANSIC)))
	row.AddCell().SetString(order.Status.String())
	row.AddCell().SetString(order.PaypallToken)
	row.AddCell().SetString(order.Comment)

	file.Save(excel.FilePath)
}

func (excel *ExcelOrderDao) ParseOrders() {
	xlFile, _ := xlsx.OpenFile(excel.FilePath)
	sheet, _ := getSheetByName(xlFile, "orders")
	excel.orders = parseOrdersSheet(sheet)
}

func parseOrdersSheet(sheet *xlsx.Sheet) []*models.Order {
	var orders []*models.Order

	for rowNum, row := range sheet.Rows {
		if rowNum == 0 {
			continue
		}

		order := new(models.Order)
		for i, cell := range row.Cells {
			cell.Value = strings.TrimSpace(cell.Value)
			switch i {
			case 0:
				order.Id = cell.Value
			case 1:
				order.ItemId = cell.Value
			case 2:
				order.Count = strconv.ParseFloat(cell.Value, 64)
			case 3:
				order.Buyer.Name = cell.Value
			case 4:
				order.Delivery.DeliveryType = order.Delivery.DeliveryType.Parse(cell.Value)
			case 5:
				order.Delivery.Address = cell.Value
			case 6:
				order.Datetime, _ = time.Parse(time.ANSIC, cell.Value)
			case 7:
				order.Status = order.Status.Parse(cell.Value)
			case 8:
				order.PaypallToken = cell.Value
			case 9:
				order.Comment = cell.Value
			}
		}

		if (order.Id != "") {
			orders = append(orders, order)
		}
	}

	return orders
}
