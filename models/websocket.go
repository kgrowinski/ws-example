package models

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketAction string

const (
	WS_INIT      WebsocketAction = "INIT_CONNECTION"
	WS_MESSAGE   WebsocketAction = "MESSAGE"
	WS_SET_COLOR WebsocketAction = "SET_COLOR"
	WS_NEW_COLOR WebsocketAction = "NEW_COLOR"
	WS_PING      WebsocketAction = "PING"
	WS_PONG      WebsocketAction = "PONG"
	WS_ERROR     WebsocketAction = "ERROR"
)

// Time allowed to read the next pong message from the peer.
var PONG_WAIT = 10 * time.Second

// Send pings to peer with this period. Must be less than pongWait.
var PING_INTERVAL = (PONG_WAIT * 9) / 10

var WebsocketManager *WSManager = &WSManager{
	Clients: make(map[*WSClient]bool),
}

type WSManager struct {
	Clients map[*WSClient]bool
	sync.RWMutex
}

type WSClient struct {
	ID         string
	Connection *websocket.Conn
	Manager    *WSManager

	Egress chan WSResponse[WSResponsePayload]
}

func NewWSClient(ws *websocket.Conn, manager *WSManager) *WSClient {
	return &WSClient{
		Connection: ws,
		Manager:    manager,

		Egress: make(chan WSResponse[WSResponsePayload]),
	}
}

func (c *WSClient) SetID(id string) {
	c.ID = id
}

func (c *WSClient) Send(v WSResponse[WSResponsePayload]) {
	c.Egress <- v
}

func (c *WSClient) HandlePong(pongMsg string) error {
	return c.Connection.SetReadDeadline(time.Now().Add(PONG_WAIT))
}

func (c *WSClient) HandleError(err error, appDomain string, errAppCode int) {
	if err != nil {
		fmt.Println(err)
	}

	message := WSResponse[WSResponsePayload]{
		Action: WS_ERROR,
		Payload: WSErrorMessage{
			ErrorCode: errAppCode,
			AppDomain: appDomain,
		},
	}
	c.Send(message)
}

func ReadJSON[T WSPayload](message WSMessage[WSPayload]) (WSMessage[T], error) {
	readMessage, err := json.Marshal(message)
	if err != nil {
		return WSMessage[T]{}, err
	}

	var newMessage WSMessage[T]

	err = json.Unmarshal(readMessage, &newMessage)
	if err != nil {
		return WSMessage[T]{}, err
	}

	return newMessage, nil

}

type WSPayload interface {
	WSPaginationMessage | WSPingMessage | string | any
}

type WSResponsePayload interface {
	WSErrorMessage | WSFileUploaded | string | any
}

type WSMessage[T WSPayload] struct {
	Action        WebsocketAction `json:"action"`
	Authorization string          `json:"authorization"`
	Payload       T               `json:"payload"`
}

type WSResponse[T WSResponsePayload] struct {
	Action  WebsocketAction `json:"action"`
	Payload T               `json:"payload"`
}

type WSErrorMessage struct {
	AppDomain string `json:"appDomain"`
	ErrorCode int    `json:"errorCode"`
}

type WSFileUploaded struct {
	Type string `json:"type"`
}

type WSPaginationMessage struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type WSPingMessage struct {
	PingID string `json:"pingId"`
}

type WSPongMessage struct {
	PongID string `json:"pongId"`
}
