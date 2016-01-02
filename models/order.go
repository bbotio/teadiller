package models

type DeliveryType int

const (
    SelfService DeliveryType = iota
    RussianPostOffice DeliveryType = iota
)

type OrderStatus int

const (
    New OrderStatus = iota
    Paid OrderStatus = iota
    Processing OrderStatus = iota
    Shipped OrderStatus = iota
    Delivered OrderStatus = iota
    Done OrderStatus = iota
    Cancelled OrderStatus = iota
)

type Buyer struct {
    Name string
}

type Delivery struct {
    DeliveryType DeliveryType
    Address string // optional 
}

type Order struct {
    Id string
    ItemId string
    Buyer Buyer
    Delivery Delivery
    Datetime uint64
    Status OrderStatus
    PaypallToken string
    Comment string // comments from seller
}
