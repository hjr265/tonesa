package ui

import (
	"html/template"

	"github.com/hjr265/tonesa/data"
)

type TplIndexValues struct{}

var TplIndex = template.Must(template.ParseFiles("ui/templates/layout.html.tpl", "ui/templates/index.html.tpl"))

type TplUploadViewValues struct {
	Upload *data.Upload
}

var TplUploadView = template.Must(template.ParseFiles("ui/templates/layout.html.tpl", "ui/templates/uploadView.html.tpl"))
