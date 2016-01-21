package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    "teadiller/botflow"
    "teadiller/web"
    "teadiller/flows"
)

func main() {

    initFlow := botflow.Flow{Command: "", Handler: flows.Default}
    initFlow.Bind("/about", flows.About)
    log.Printf("Bot Flow: %s", initFlow)


    done := make(chan bool)
    err := botflow.StartBot(os.Getenv("TELEGRAMM_TOKEN"), os.Getenv("TELEGRAM_BOT_NAME"), initFlow, done)
    if err != nil {
        panic(err)
    }

    web.StartServer(os.Getenv("BOT_WEB_PORT"))

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
