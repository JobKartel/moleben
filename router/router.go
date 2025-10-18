package router

import (
	"moleben/controller"

	"github.com/gin-gonic/gin"
)

type Router struct{ Engine *gin.Engine }

func New(ctrl *controller.ChatController) *Router {
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	r.GET("/health", func(c *gin.Context) { c.String(200, "ok") })

	api := r.Group("/api")
	{
		api.POST("/sessions", ctrl.CreateSession)
		sr := api.Group("/sessions/:id")
		{
			sr.GET("/", ctrl.GetSession)
			sr.POST("/messages", ctrl.PostMessage)
			sr.GET("/messages", ctrl.ListMessages)
		}
	}
	return &Router{Engine: r}
}
