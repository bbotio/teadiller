package models

type ItemType int

const (
    CountItem ItemType = iota
    VolumeItem ItemType = iota
)

type Item struct {
    Id string
    Name string
    Desc string
    PhotoPath string
    Tags []string
    Type ItemType
    Count float64
    AdditionalFields map[string]string
}

type ItemDao interface {
	GetById(id string) (*Item, error)
	GetAll() []*Item
}