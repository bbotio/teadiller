package flows

import (
	"teadiller/Godeps/_workspace/src/github.com/Syfaro/telegram-bot-api"
    "teadiller/Godeps/_workspace/src/github.com/jasonwinn/geocoder"
    "teadiller/botflow"
    "teadiller/payment"
    "teadiller/models"
	"errors"
	"fmt"
	"os"
	"log"
)

var EMPTY_BODY = []tgbotapi.Chattable{}

func Location(msg tgbotapi.Message, ctx botflow.Context) ([]tgbotapi.Chattable, error) {
	order := ctx["order"].(models.Order)

	if msg.Location.Latitude == 0 && msg.Location.Longitude == 0 {
		_, _, err := geocoder.Geocode(msg.Text)
		order.Delivery.Address = msg.Text

		if err != nil {
			log.Printf("Can't validate address: %s", msg.Text)
		}
	} else {
		location, err := geocoder.ReverseGeocode(float64(msg.Location.Latitude), float64(msg.Location.Longitude))

		if err == nil {
			order.Delivery.Address = fmt.Sprintf("%+v", location)
		} else {
			return EMPTY_BODY, err
		}
	}

	item := ctx["item"].(models.Item)
	url, err := preparePaymentUrl(item)

	if err != nil {
		return EMPTY_BODY, err
	}

	orderDao := models.GetOrderDao()
	order.PaypalToken = url[len(url) - 20:]
	orderDao.Save(&order)

	paymentUrl := "http://localhost:" + os.Getenv("TELEGRAM_BOT_WEB_PORT") + "/buy/" + order.Id
	return []tgbotapi.Chattable{
		tgbotapi.NewMessage(msg.Chat.ID, order.Delivery.Address),
		tgbotapi.NewMessage(msg.Chat.ID, "Payment url: " + paymentUrl),
	}, nil
}

func preparePaymentUrl(item models.Item) (url string, err error) {
	response, err := payment.SetExpressCheckout(item)

	if err != nil {
		return "Error on setExpressCheckout request", err
	}

	if response.Ack == "Success" {
		return response.CheckoutUrl(), nil
	}

	return "", errors.New("Response Ack != Success, " + response.Ack)
}
