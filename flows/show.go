package flows

import (
    "teadiller/botflow"
    "teadiller/Godeps/_workspace/src/github.com/Syfaro/telegram-bot-api"
    "teadiller/models"
)


func Show(msg tgbotapi.Message, ctx botflow.Context) ([]tgbotapi.MessageConfig, error) {
    itemDao := models.GetItemDao()
    if (msg.ReplyToMessage != nil) {
        categoryName := msg.ReplyToMessage.Text
        items := itemDao.GetByCategory(categoryName)

        if len(items) > 0 {
            result := []tgbotapi.MessageConfig{}
            for _, item := range items {
                result = append(result, tgbotapi.NewMessage(msg.Chat.ID, item.Name + " http://localhost:8080/items/" + item.Id))
            }
            return result, nil
        }
    }

    return []tgbotapi.MessageConfig{tgbotapi.NewMessage(msg.Chat.ID, "Sorry I don't know what to show")}, nil
}
