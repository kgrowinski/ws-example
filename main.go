package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"websockets.com/routes"
)

type App struct {
	Router *gin.Engine
}

func (a *App) CreateRoutes(cache *cache.Cache) {
	a.Router = gin.New()
	a.Router.LoadHTMLGlob("templates/*.html")

	direct := a.Router.Group("/")
	direct.GET("", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title":   "Page 1 IS HERE",
			"content": "page1.html",
		})
	})
	ws := a.Router.Group("/ws/v1")
	routes.CreateWebsocketRoutes(ws, cache)
}

func (a *App) Run() {
	port := "8080"

	cache := cache.New(60*time.Minute, 5*time.Minute)

	a.CreateRoutes(cache)

	a.Router.Run(fmt.Sprintf(":%s", port))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	app := App{}
	app.Run()
}
