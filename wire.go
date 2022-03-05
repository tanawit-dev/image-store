//+build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/tanawit-dev/image-store/common/db"
	"github.com/tanawit-dev/image-store/common/minio"
	"github.com/tanawit-dev/image-store/config"
	"github.com/tanawit-dev/image-store/repositories"
	"github.com/tanawit-dev/image-store/routers"
	"github.com/tanawit-dev/image-store/server"
	"github.com/tanawit-dev/image-store/services"
)

func InitializeApplication() (*gin.Engine, error) {
	wire.Build(
		config.ProvideConfig,
		db.ProvideDatabase,
		minio.ProvideImageStorage,
		repositories.ProvideImageRepository,
		services.ProvideImageService,
		routers.ProvideImageController,
		server.ProvideServer,
	)

	return &gin.Engine{}, nil
}
