package flows
import (
	"teadiller/Godeps/_workspace/src/github.com/Syfaro/telegram-bot-api"
	"teadiller/botflow"
	"teadiller/models"
	"errors"
	"fmt"
)

func DeliveryType(msg tgbotapi.Message, ctx botflow.Context) ([]tgbotapi.Chattable, error) {
	order := ctx["order"].(models.Order)
	switch msg.Text {
	case "Russian Post":
		order.Delivery.DeliveryType = models.RussianPostOffice
	case "Self-service":
		order.Delivery.DeliveryType = models.SelfService
	default:
		return EMPTY_BODY,
		errors.New(fmt.Sprintf("Invalid delivery type, %s", msg.Text))
	}

	ctx["order"] = order
	return []tgbotapi.Chattable{tgbotapi.NewMessage(msg.Chat.ID, "Please, enter delivery address")}, nil
}
