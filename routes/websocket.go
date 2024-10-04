package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"

	"websockets.com/controllers"
)

func CreateWebsocketRoutes(routerGroup *gin.RouterGroup, cache *cache.Cache) {

	WebsocketController := controllers.NewWebsocketController(cache)

	WebsocketRouter := routerGroup.Group("/websocket")
	{
		WebsocketRouter.GET("", WebsocketController.GetWebsocket)
	}
}
