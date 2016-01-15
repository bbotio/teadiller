package web

import (
    "fmt"
    "net/http"
    "teadiller/web/controllers"
)

func StartServer(port string) {
    http.HandleFunc("/", controllers.Welcome)
    go func() {
        http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
    }()
}
