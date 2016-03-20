package controllers

import (
    "log"
    "net/http"
    "html/template"
    "teadiller/web/templates"
    "teadiller/Godeps/_workspace/src/github.com/gorilla/mux"
    "teadiller/models"
)

func Item(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    itemId,_ := vars["itemId"]
    log.Printf("Got item with id %s", itemId)
    itemDao := models.GetItemDao()
    item, err := itemDao.GetById(itemId)
    if err != nil {
        errorTmpl, _ := template.New("error").Parse(templates.Error)
        errorTmpl.Execute(w, nil)
        return
    }
    t, _ := template.New("item").Parse(templates.Item)
    t.Execute(w, *item)
}
