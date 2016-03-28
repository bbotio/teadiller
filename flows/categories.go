package flows

import (
    "teadiller/botflow"
    "teadiller/Godeps/_workspace/src/github.com/deckarep/golang-set"
    "teadiller/Godeps/_workspace/src/github.com/Syfaro/telegram-bot-api"
    "teadiller/models"
)


func Categories(msg tgbotapi.Message, ctx botflow.Context) ([]tgbotapi.Chattable, error) {
    itemDao := models.GetItemDao()
    items := itemDao.GetAll()
    categoriesMessages := mapset.NewSet()

    for _, item := range items {
        for _, tag := range item.Tags {
            categoriesMessages.Add(tag)
        }
    }

    result := []tgbotapi.Chattable{}
    for category := range categoriesMessages.Iter() {
        msg := tgbotapi.NewMessage(msg.Chat.ID, category.(string))
        result = append(result, msg)
    }

    return result, nil
}
