package templates

import (
	"html/template"
)

var Main = template.Must(template.ParseFiles("www/templates/index.html"))
