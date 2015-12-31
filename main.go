package main

import (
	"fmt"
	"net/http"
	"os"
	"teadiller/Godeps/_workspace/src/github.com/parnurzeal/gorequest"
)

const TELEGRAMM_URL_PATTERN = "https://api.telegram.org/bot%s/"
const BOT_URL = ""

func main() {
	request := gorequest.New()
	fmt.Println("register bot...")
	url := fmt.Sprintf(TELEGRAMM_URL_PATTERN, os.Getenv("TELEGRAMM_TOKEN"))
	resp, _, _ := request.Get(fmt.Sprint("%ssetWebhook?url=%s", url, BOT_URL)).End()
	if resp.Status == "200" {
		fmt.Println("listening...")
		err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
		if err != nil {
			panic(err)
		}
	}
}

func hook(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Got it")
}
