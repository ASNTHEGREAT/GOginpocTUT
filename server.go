package main

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ASNTHEGREAT/golang-gin-poc/controller"
	"github.com/ASNTHEGREAT/golang-gin-poc/middleware"
	"github.com/ASNTHEGREAT/golang-gin-poc/service"
	gindump "github.com/tpkeeper/gin-dump"
)

var(
	videoService service.VideoService = service.New()
	videoController controller.VideoController = controller.New(videoService)
) 

func setupLogOutput() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main() {

	setupLogOutput()

	server := gin.New()

	server.Static("/css", "./templates/css")

	server.LoadHTMLGlob("templates/*.html")
	
	server.Use(gin.Recovery(), middleware.Logger(),
		middleware.BasicAuth(), gindump.Dump())


	apiRoutes := server.Group("/api") 
	{
		apiRoutes.GET("/posts", func (ctx *gin.Context)  {
			ctx.IndentedJSON(http.StatusOK, videoController.FindAll())
		})
	
		apiRoutes.POST("/videos", func (ctx *gin.Context)  {
			err := videoController.Save(ctx)
			if err != nil {	
				ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Video Input is Valid!!"})
			}
	
		})
	}

	viewRoutes := server.Group("/view")
	{
		viewRoutes.GET("/videos", videoController.ShowAll)
	}

	server.Run(":8000")
}