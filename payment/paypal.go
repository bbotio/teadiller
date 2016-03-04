package payment

import (
	"teadiller/Godeps/_workspace/src/github.com/crowdmob/paypal"
    "teadiller/Godeps/_workspace/src/github.com/gorilla/mux"
    "teadiller/models"
    "sync"
    "net/http"
	"fmt"
	"log"
)

var client *paypal.PayPalClient
var onceForItem sync.Once
var okUrl string
var cancelUrl string

var redirectUrlTemplate = "https://www.sandbox.paypal.com/cgi-bin/webscr?cmd=_express-checkout&token=%s"

func InitClient(username string, password string, signature string, sandbox bool, ok string, cancel string) *paypal.PayPalClient {
	onceForItem.Do(func() {
		client = paypal.NewDefaultClient(username, password, signature, sandbox)
	})

	okUrl = ok
	cancelUrl = cancel

	return client
}

func SetExpressCheckout(item models.Item) (*paypal.PayPalResponse, error) {
	currentItem := []paypal.PayPalDigitalGood{paypal.PayPalDigitalGood{
		Name: item.Name,
		Amount: item.Price,
		Quantity: int16(item.Count),
	}}

	return client.SetExpressCheckoutDigitalGoods(paypal.SumPayPalDigitalGoodAmounts(&currentItem),
		"USD",
		okUrl,
		cancelUrl,
		currentItem,
	)
}

func DoExpressCheckout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	//orderDao := models.GetOrderDao()
	//order, err := orderDao.GetByPayPalToken(token)

	response, err := client.DoExpressCheckoutSale(token, vars["payerId"], "USD", 1) // fix it!!

	if err == nil {
		log.Printf("SUCCESS %+v", response)
	} else {
		log.Printf("ERROR %+v %+v", response, err)
	}
}

func CancelledCheckout(w http.ResponseWriter, r *http.Request) {

}

func PayPalRedirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderId, _ := vars["orderId"]

	orderDao := models.GetOrderDao()
	order, err := orderDao.GetById(orderId)

	log.Printf("%+v%+v", order, err)

	if err == nil {
		http.Redirect(w, r, fmt.Sprintf(redirectUrlTemplate, order.PaypalToken), 301)
	}
}

