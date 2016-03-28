package flows

import (
    "teadiller/botflow"
    "teadiller/Godeps/_workspace/src/github.com/Syfaro/telegram-bot-api"
    "teadiller/models"
)


func Show(msg tgbotapi.Message, ctx botflow.Context) ([]tgbotapi.Chattable, error) {
    itemDao := models.GetItemDao()
    if (msg.ReplyToMessage != nil) {
        categoryName := msg.ReplyToMessage.Text
        items := itemDao.GetByCategory(categoryName)

        if len(items) > 0 {
            result := []tgbotapi.Chattable{}
            for _, item := range items {
                pc := tgbotapi.NewPhotoUpload(msg.Chat.ID, item.PhotoPath)
                pc.Caption = "#" + item.Id + " - " + item.Name + "\n" + item.Desc
                result = append(result, pc)
            }
            return result, nil
        }
    }

    return []tgbotapi.Chattable{tgbotapi.NewMessage(msg.Chat.ID, "Sorry I don't know what to show")}, nil
}
