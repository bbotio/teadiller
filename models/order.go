package models

import (
	"time"
	"path/filepath"
	"github.com/tealeg/xlsx"
	"fmt"
	"errors"
)

type DeliveryType string
type OrderStatus string

const (
	SelfService DeliveryType = "Self-service"
	RussianPostOffice DeliveryType = "Russian Post Office"
)

const (
	New OrderStatus = "New"
	Paid OrderStatus = "Paid"
	Processing OrderStatus = "Processing"
	Shipped OrderStatus = "Shipped"
	Delivered OrderStatus = "Delivered"
	Done OrderStatus = "Done"
	Cancelled OrderStatus = "Cancelled"
)

type Buyer struct {
	Name string
}

type Delivery struct {
	DeliveryType DeliveryType
	Address      string // optional
}

type Order struct {
	Id           string
	ItemId       string
	Buyer        Buyer
	Delivery     Delivery
	Datetime     time.Time
	Status       OrderStatus
	PaypallToken string
	Comment      string // comments from seller
}

func (p DeliveryType) Parse(deliveryType string) DeliveryType {
	switch deliveryType {
	case "Self-service":
		return SelfService
	case "Russian Post Office":
		return RussianPostOffice
	}
	return SelfService
}

func (p DeliveryType) String() string {
	return string(p)
}

func (p OrderStatus) String() string {
	return string(p)
}

func (p OrderStatus) Parse(code string) OrderStatus {
	switch code {
	case "New":
		return New
	case "Paid":
		return Paid
	case "Processing":
		return Processing
	case "Shipped":
		return Shipped
	case "Delivered":
		return Delivered
	case "Done":
		return Done
	case "Cancelled":
		return Cancelled
	}

	return New
}

func (order Order) Save(filepath string) {
	file, _ := xlsx.OpenFile(filepath)

	sheet, _ := getSheetByName(file, "orders")
	row := sheet.AddRow();

	row.AddCell().SetString(order.Id)
	row.AddCell().SetString(order.ItemId)
	row.AddCell().SetString(order.Buyer.Name)
	row.AddCell().SetString(order.Delivery.DeliveryType.String())
	row.AddCell().SetString(order.Delivery.Address)
	row.AddCell().SetString(fmt.Sprint(order.Datetime.Format(time.ANSIC)))
	row.AddCell().SetString(order.Status.String())
	row.AddCell().SetString(order.PaypallToken)
	row.AddCell().SetString(order.Comment)

	file.Save(filepath)
}

func getSheetByName(file *xlsx.File, sheetName string) (*xlsx.Sheet, error) {
	for _, sheet := range file.Sheets {
		if sheet.Name == "orders" {
			return sheet, nil
		}
	}

	return nil, errors.New("Sheet with name = " + sheetName + " not found.")
}
