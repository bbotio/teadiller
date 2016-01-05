package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    "teadiller/Godeps/_workspace/src/github.com/Syfaro/telegram-bot-api"
//    "teadiller/models"
      "teadiller/botflow"
	"log"
	"teadiller/Godeps/_workspace/src/github.com/Syfaro/telegram-bot-api"
	"teadiller/Godeps/_workspace/src/github.com/tealeg/xlsx"
	"teadiller/models"
	"strings"
	"strconv"
	"fmt"
	"time"
	"errors"
)

func main() {
    handler := func(msg tgbotapi.Message, ctx botflow.Context) (tgbotapi.MessageConfig, error) {
                       return tgbotapi.NewMessage(msg.Chat.ID, "Would you like tea?"), nil
                    }

    initFlow := botflow.Flow{Command: "", Handler: handler}
    aboutHandler := func(msg tgbotapi.Message, ctx botflow.Context) (tgbotapi.MessageConfig, error) {
                       return tgbotapi.NewMessage(msg.Chat.ID, "I'm tea bot I wanna sale you everything"), nil
                    }
    initFlow.Bind("/about", aboutHandler)

    log.Printf("Bot Flow: %s", initFlow)
    done := make(chan bool)
    err := botflow.StartBot(os.Getenv("TELEGRAMM_TOKEN"), os.Getenv("TELEGRAM_BOT_NAME"), initFlow, done)
    if err != nil {
        panic(err)
    }

    signalChannel := make(chan os.Signal, 2)
    signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

    sig := <-signalChannel
    switch sig {
        case os.Interrupt:
            log.Printf("Got Interupt")
        case syscall.SIGTERM:
            log.Printf("Got SIGTERM")
    }
    log.Printf("Stop bot")
    done <- true
    log.Printf("Bye! Bye!")
}

		orders[0].Save("asd.xlsx")
	}
}

// return list of items and orders
func parseXslx(filepath string) ([]models.Item, []models.Order) {
	xlFile, err := xlsx.OpenFile(filepath)

	if err != nil {
	}

	var items []models.Item
	var orders []models.Order

	for _, sheet := range xlFile.Sheets {
		if sheet.Name == "items" {
			items = parseItems(sheet)
		} else if sheet.Name == "orders" {
			orders = parseOrders(sheet)
		}
	}

	return items, orders
}

func parseItems(sheet *xlsx.Sheet) []models.Item {
	var items []models.Item

	for rowNum, row := range sheet.Rows {
		if rowNum == 0 {
			continue
		}

		var item models.Item
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
				item.Type = models.VolumeItem
			case 6:
				item.AdditionalFields = make(map[string]string)
				properties := strings.Fields(cell.Value)
				for _, property := range properties {
					if len(property) == 0 {
						continue
					}

					keyvalue := strings.Split(property, "=")
					item.AdditionalFields[keyvalue[0]] = keyvalue[1]
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

func parseOrders(sheet *xlsx.Sheet) []models.Order {
	var orders []models.Order

	for rowNum, row := range sheet.Rows {
		if rowNum == 0 {
			continue
		}

		var order models.Order
		for i, cell := range row.Cells {
			cell.Value = strings.TrimSpace(cell.Value)
			switch i {
			case 0:
				order.Id = cell.Value
			case 1:
				order.ItemId = cell.Value
			case 2:
				order.Buyer.Name = cell.Value
			case 3:
				order.Delivery.DeliveryType = order.Delivery.DeliveryType.Parse(cell.Value)
			case 4:
				order.Delivery.Address = cell.Value
			case 5:
				order.Datetime, _ = time.Parse(time.ANSIC, cell.Value)
			case 6:
				order.Status = order.Status.Parse(cell.Value)
			case 7:
				order.PaypallToken = cell.Value
			case 8:
				order.Comment = cell.Value
			}
		}

		if (order.Id != "") {
			orders = append(orders, order)
		}
	}

	return orders
}

func getItemById(id string) (models.Item, error) {
	for _, item := range items {
		if item.Id == id {
			return item, nil
		}
	}

	var result models.Item
	return result, errors.New("Item with id = " + id + " not found.")
}

func getOrderById(id string) (models.Order, error) {
	for _, order := range orders {
		if order.Id == id {
			return order, nil
		}
	}
	var result models.Order
	return result, errors.New("Order with id = " + id + " not found.")
}