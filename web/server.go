package web

import (
    "log"
    "fmt"
    "net/http"
    "teadiller/web/controllers"
    "teadiller/payment"
    "teadiller/Godeps/_workspace/src/github.com/gorilla/mux"
)

func StartServer(port string) {
    r := mux.NewRouter()

    r.HandleFunc("/", controllers.Welcome)
	r.HandleFunc("/buy/{orderId}", payment.PayPalRedirectHandler)
	r.HandleFunc("/ok?token={token}&PayerID={payerId}", payment.DoExpressCheckout)
	r.HandleFunc("/cancel?token={token}", payment.CancelledCheckout)
    r.HandleFunc("/items/{itemId}", controllers.Item)
    go func() {
        log.Printf("Listen on port %s ...", port)
        http.ListenAndServe(fmt.Sprintf(":%s", port), r)
    }()
}
