package web

import (
    "teadiller/Godeps/_workspace/src/github.com/gorilla/mux"
    "teadiller/web/controllers"
    "teadiller/payment"
    "net/http"
    "fmt"
)

func StartServer(port string) {
	r := mux.NewRouter()
	r.HandleFunc("/", controllers.Welcome)
	r.HandleFunc("/buy/{orderId}", payment.PayPalRedirectHandler)
	r.HandleFunc("/ok?token={token}&PayerID={payerId}", payment.DoExpressCheckout)
	r.HandleFunc("/cancel?token={token}", payment.CancelledCheckout)

	go func() {
        http.ListenAndServe(fmt.Sprintf(":%s", port), r)
    }()
}
