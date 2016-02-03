package flows

import (
    "teadiller/botflow"
//    "teadiller/Godeps/_workspace/src/github.com/deckarep/golang-set"
    "teadiller/Godeps/_workspace/src/github.com/Syfaro/telegram-bot-api"
//    "teadiller/models"
)


func Show(msg tgbotapi.Message, ctx botflow.Context) ([]tgbotapi.MessageConfig, error) {
    return []tgbotapi.MessageConfig{tgbotapi.NewMessage(msg.Chat.ID, "Some shit")}, nil
}
