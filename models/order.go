package models

import (
	"time"
    "sync"
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

func (dt DeliveryType) Parse(deliveryType string) DeliveryType {
	switch deliveryType {
	case "Self-service":
		return SelfService
	case "Russian Post Office":
		return RussianPostOffice
	}
	return SelfService
}

func (deliveryType DeliveryType) String() string {
	return string(deliveryType)
}

func (orderStatus OrderStatus) String() string {
	return string(orderStatus)
}

func (os OrderStatus) Parse(code string) OrderStatus {
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

type OrderDao interface {
	GetById(id string) (*Order, error)
	GetAll() []*Order
	Save(order *Order)
}


var orderDao OrderDao
var onceForOrder sync.Once

func InitOrderDao(initOrderDao OrderDao) OrderDao{
    onceForOrder.Do(func(){
            orderDao = initOrderDao
        })
    return orderDao
}

func GetOrderDao() OrderDao {
    return orderDao
}
