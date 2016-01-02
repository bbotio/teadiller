package main

import (
    "log"
    "os"
    "teadiller/Godeps/_workspace/src/github.com/Syfaro/telegram-bot-api"
//    "teadiller/models"
)

func main() {
    bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAMM_TOKEN"))
    if err != nil {
        panic(err)
    }

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates, err := bot.GetUpdatesChan(u)

    for update := range updates {
        log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
    }
}
