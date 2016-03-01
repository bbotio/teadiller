package models

import(
    "sync"
)

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
	Price float64
    AdditionalFields map[string]string
}

type ItemDao interface {
	GetById(id string) (*Item, error)
	GetAll() []*Item
}

var itemDao ItemDao
var onceForItem sync.Once

func InitItemDao(initItemDao ItemDao) ItemDao{
    onceForItem.Do(func(){
            itemDao = initItemDao
        })
    return itemDao
}

func GetItemDao() ItemDao {
    return itemDao
}
