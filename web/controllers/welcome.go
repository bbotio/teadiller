package controllers

import (
    "net/http"
    "html/template"
    "teadiller/web/templates"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
    t, _ := template.New("welcome").Parse(templates.Welcome)
    t.Execute(w, struct{Text string}{Text: "Welcome to bot kingdom"})
}
