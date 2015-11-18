package ui

import (
	"net/http"
)

var AssetsFS = http.FileServer(http.Dir("ui/assets"))
