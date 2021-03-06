// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from wire.go:

func InitializeApplication() (*gin.Engine, error) {
	configuration := config.ProvideConfig()
	gormDB, err := db.ProvideDatabase(configuration)
	if err != nil {
		return nil, err
	}
	imageRepository := repositories.ImageRepository{
		DB: gormDB,
	}
	imageStorage, err := minio.ProvideImageStorage(configuration)
	if err != nil {
		return nil, err
	}
	imageService := services.ImageService{
		Repo:         imageRepository,
		ImageStorage: imageStorage,
		DB:           gormDB,
	}
	imageController := routers.ImageController{
		Service: imageService,
	}
	engine, err := server.ProvideServer(imageController)
	if err != nil {
		return nil, err
	}
	return engine, nil
}

// wire.go:

var ImageRepositoryStructProvider = wire.Struct(new(repositories.ImageRepository), "*")

var ImageServiceStructProvider = wire.Struct(new(services.ImageService), "*")

var ImageControllerStructProvider = wire.Struct(new(routers.ImageController), "*")

var ImageModuleSet = wire.NewSet(ImageControllerStructProvider, ImageServiceStructProvider, ImageRepositoryStructProvider)

var InfraSet = wire.NewSet(config.ProvideConfig, db.ProvideDatabase, minio.ProvideImageStorage, server.ProvideServer)

var ModuleSet = wire.NewSet(ImageModuleSet)

var RootSet = wire.NewSet(InfraSet, ModuleSet)
