package templates

import (
	"html/template"
	"net/http"
)

var parsedTemplates = template.Must(template.ParseFiles(
	"templates/home.html",
	"templates/register.html",
	"templates/login.html",
	"templates/users.html"))

type PageData struct {
	Data interface{}
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := parsedTemplates.ExecuteTemplate(w, tmpl+".html", PageData{
		Data: data,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
