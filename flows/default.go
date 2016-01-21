package flows

import (
    "teadiller/botflow"
    "teadiller/Godeps/_workspace/src/github.com/Syfaro/telegram-bot-api"
)


func Default(msg tgbotapi.Message, ctx botflow.Context) ([]tgbotapi.MessageConfig, error) {
    return []tgbotapi.MessageConfig{tgbotapi.NewMessage(msg.Chat.ID, "Would you like tea?")}, nil
}
