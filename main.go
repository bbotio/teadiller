package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    "teadiller/Godeps/_workspace/src/github.com/Syfaro/telegram-bot-api"
//    "teadiller/models"
    "teadiller/botflow"
)

func main() {
    handler := func(msg tgbotapi.Message, ctx botflow.Context) ([]tgbotapi.MessageConfig, error) {
                       return []tgbotapi.MessageConfig{tgbotapi.NewMessage(msg.Chat.ID, "Would you like tea?")}, nil
                    }

    initFlow := botflow.Flow{Command: "", Handler: handler}
    aboutHandler := func(msg tgbotapi.Message, ctx botflow.Context) ([]tgbotapi.MessageConfig, error) {
                       return []tgbotapi.MessageConfig{tgbotapi.NewMessage(msg.Chat.ID, "I'm tea bot I wanna sale you everything"),
                                                        tgbotapi.NewMessage(msg.Chat.ID, "Really")}, nil
                    }
    initFlow.Bind("/about", aboutHandler)

    log.Printf("Bot Flow: %s", initFlow)
    done := make(chan bool)
    err := botflow.StartBot(os.Getenv("TELEGRAMM_TOKEN"), os.Getenv("TELEGRAM_BOT_NAME"), initFlow, done)
    if err != nil {
        panic(err)
    }

    signalChannel := make(chan os.Signal, 2)
    signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

    sig := <-signalChannel
    switch sig {
        case os.Interrupt:
            log.Printf("Got Interupt")
        case syscall.SIGTERM:
            log.Printf("Got SIGTERM")
    }
    log.Printf("Stop bot")
    done <- true
    log.Printf("Bye! Bye!")
}