package flows

import (
	"teadiller/botflow"
	"teadiller/Godeps/_workspace/src/github.com/Syfaro/telegram-bot-api"
	"fmt"
	"teadiller/models"
	"strconv"
	"errors"
)

var services = [][]string{{"Russian Post", "Self-service"}}
var keyboard = tgbotapi.ReplyKeyboardMarkup{
	Keyboard:services,
	ResizeKeyboard:true,
	OneTimeKeyboard:true,
	Selective:true,
}

func Count(msg tgbotapi.Message, ctx botflow.Context) ([]tgbotapi.Chattable, error) {
	if _, ok := ctx["item"]; !ok {
		return EMPTY_BODY, errors.New("Please, choose item")
	}

	item := ctx["item"].(models.Item)
	order := ctx["order"].(models.Order)
	requestedCount, err := strconv.ParseFloat(msg.Text, 64)

	if err != nil {
		return EMPTY_BODY,
		errors.New(fmt.Sprintf("Error parsing requested count = %s. Please, enter valid number!", msg.Text))
	}

	// TODO: Think about better error value
	EPS := 1e-2

	if requestedCount >= item.Count + EPS || requestedCount + EPS <= 0.0 {
		return EMPTY_BODY,
		errors.New(fmt.Sprintf("You've entered invalid count, requested count = %d, available count %f",
			requestedCount, item.Count))
	}

	order.Count = requestedCount
	ctx["order"] = order

	message := tgbotapi.NewMessage(msg.Chat.ID, "Please, choose delivery type")
	message.ReplyMarkup = keyboard
	return []tgbotapi.Chattable{message}, nil
}

