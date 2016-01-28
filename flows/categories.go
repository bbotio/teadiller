package flows

import (
    "teadiller/botflow"
    "teadiller/Godeps/_workspace/src/github.com/deckarep/golang-set"
    "teadiller/Godeps/_workspace/src/github.com/Syfaro/telegram-bot-api"
    "teadiller/models"
)


func Categories(msg tgbotapi.Message, ctx botflow.Context) ([]tgbotapi.MessageConfig, error) {
    itemDao := models.GetItemDao()
    items := itemDao.GetAll()
    categorieMessages := mapset.NewSet()

    for _, item := range items {
        for _, tag := range item.Tags {
            categorieMessages.Add(tag)
        }
    }

    result := []tgbotapi.MessageConfig{}
    for category := range categorieMessages.Iter() {
        msg := tgbotapi.NewMessage(msg.Chat.ID, category.(string))
        result = append(result, msg)
    }

    return result, nil
}
