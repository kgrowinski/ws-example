package helpers

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var DATE_FORMAT = "01/02/2006"

var WS_UPGRADER = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
