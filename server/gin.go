package server

import (
	"github.com/gin-gonic/gin"
	"github.com/tanawit-dev/image-store/routers"
)

func ProvideServer(imageController routers.ImageController) (*gin.Engine, error) {
	r := gin.Default()

	v1 := r.Group("/api")
	registerImageController(v1.Group("/images"), imageController)

	return r, nil
}

func registerImageController(group *gin.RouterGroup, controller routers.ImageController) {
	group.GET("/:id", controller.GetImage)
	group.POST("/", controller.UploadImage)
}
