package flows

import (
	"teadiller/botflow"
	"teadiller/Godeps/_workspace/src/github.com/Syfaro/telegram-bot-api"
	"fmt"
	"errors"
	"teadiller/models"
	"strconv"
	"time"
)

func Buy(msg tgbotapi.Message, ctx botflow.Context) ([]tgbotapi.Chattable, error) {
	itemName := msg.ReplyToMessage.Text;

	if !isValidItem(itemName) {
		return []tgbotapi.Chattable{}, errors.New("Invalid item " + itemName)
	}

	// todo: get item by name
	item := models.Item{Id: "1", Name: itemName, Count: 123.0, Price: 15.0, Type: models.CountItem};

	order := models.Order{
		Id: fmt.Sprintf("%s%d%d", msg.From.UserName, msg.Chat.ID, time.Now().Second()),
		ItemId: item.Id,
		Buyer: models.Buyer{
			msg.From.UserName,
		},
		Status: models.New,
	}

	ctx["item"] = item
	ctx["order"] = order

	return []tgbotapi.Chattable{
		tgbotapi.NewMessage(msg.Chat.ID, "How many? Available count " + strconv.FormatFloat(item.Count, 'f', 6, 64))},
	nil
}

func isValidItem(itemName string) bool {
	// TODO: try to get item from DB by name
	return true
}
