package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tanawit-dev/image-store/common/db"
	"github.com/tanawit-dev/image-store/common/minio"
	"github.com/tanawit-dev/image-store/config"
	"github.com/tanawit-dev/image-store/routers"
	"log"
)

func main() {
	config.Init()
	db.Init()
	minio.Init()
	startServer()
}

func startServer() {
	r := gin.Default()

	v1 := r.Group("/api")

	routers.RegisterImageRouter(v1.Group("/images"))

	err := r.Run("localhost:8080")
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Printf("server is running")
}
