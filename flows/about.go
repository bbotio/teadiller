package flows

import (
    "teadiller/botflow"
    "teadiller/Godeps/_workspace/src/github.com/Syfaro/telegram-bot-api"
)

func About(msg tgbotapi.Message, ctx botflow.Context) ([]tgbotapi.MessageConfig, error) {
    return []tgbotapi.MessageConfig{tgbotapi.NewMessage(msg.Chat.ID, "I'm tea bot I wanna sale you everything"),
                                        tgbotapi.NewMessage(msg.Chat.ID, "Really")}, nil
}
