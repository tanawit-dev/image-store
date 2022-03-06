//go:build wireinject
// +build wireinject

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

var ImageRepositoryStructProvider = wire.Struct(new(repositories.ImageRepository), "*")
var ImageServiceStructProvider = wire.Struct(new(services.ImageService), "*")
var ImageControllerStructProvider = wire.Struct(new(routers.ImageController), "*")
var ImageModuleSet = wire.NewSet(ImageControllerStructProvider, ImageServiceStructProvider, ImageRepositoryStructProvider)

var InfraSet = wire.NewSet(config.ProvideConfig, db.ProvideDatabase, minio.ProvideImageStorage, server.ProvideServer)
var ModuleSet = wire.NewSet(ImageModuleSet)
var RootSet = wire.NewSet(InfraSet, ModuleSet)

func InitializeApplication() (*gin.Engine, error) {
	wire.Build(RootSet)

	return &gin.Engine{}, nil
}
