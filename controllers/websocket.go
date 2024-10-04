package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"

	"websockets.com/helpers"
	"websockets.com/models"
)

type WebsocketController struct {
	Manager *models.WSManager
	Cache   *cache.Cache
}

func NewWebsocketController(cache *cache.Cache) *WebsocketController {
	return &WebsocketController{Cache: cache, Manager: models.WebsocketManager}
}

func (controller *WebsocketController) AddClient(client *models.WSClient) {
	controller.Manager.Lock()
	defer controller.Manager.Unlock()

	controller.Manager.Clients[client] = true

}

func (controller *WebsocketController) RemoveClient(client *models.WSClient) {
	controller.Manager.Lock()
	defer controller.Manager.Unlock()

	if _, ok := controller.Manager.Clients[client]; ok {
		client.Connection.Close()
		delete(controller.Manager.Clients, client)
	}
}

func (controller *WebsocketController) GetWebsocket(c *gin.Context) {
	conn, err := helpers.WS_UPGRADER.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %", err)
		return
	}

	client := models.NewWSClient(conn, controller.Manager)

	controller.AddClient(client)

	go controller.ReadMessages(client)
	go controller.WriteMessages(client)

	fmt.Println("Active clients: ", len(controller.Manager.Clients))
}

func (controller *WebsocketController) ReadMessages(client *models.WSClient) {

	defer func() {
		log.Println("Closing connection read messages")
		controller.RemoveClient(client)
	}()

	if err := client.Connection.SetReadDeadline(time.Now().Add(models.PONG_WAIT)); err != nil {
		log.Println(err)
		return
	}

	client.Connection.SetPongHandler(client.HandlePong)

	for {
		var newMessage models.WSMessage[models.WSPayload]

		_, payload, err := client.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			break
		}

		err = json.Unmarshal(payload, &newMessage)
		if err != nil {
			log.Printf("Error while decoding payload: %v", err)
			return
		}

		switch newMessage.Action {
		case models.WS_INIT:
			{
				client.SetID(newMessage.Authorization)
				client.Send(models.WSResponse[models.WSResponsePayload]{
					Action:  models.WS_MESSAGE,
					Payload: "OK",
				})
			}

		case models.WS_SET_COLOR:
			{
				for client := range controller.Manager.Clients {
					client.Send(models.WSResponse[models.WSResponsePayload]{
						Action:  models.WS_NEW_COLOR,
						Payload: newMessage.Payload,
					})
				}
			}

		default:
			{
				client.HandleError(fmt.Errorf("no such action"), "WS", 404)
			}
		}
	}
}

func (controller *WebsocketController) WriteMessages(client *models.WSClient) {
	defer func() {
		log.Println("Closing connection write messages")
		controller.RemoveClient(client)
	}()

	ticker := time.NewTicker(models.PING_INTERVAL)

	for {
		select {
		case message, ok := <-client.Egress:
			if !ok {
				if err := client.Connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println(err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}

			if err := client.Connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
				return
			}

		case <-ticker.C:
			if err := client.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
				return
			}
		}

	}
}
