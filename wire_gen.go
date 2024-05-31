// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"LinkMe/internal/api"
	"LinkMe/internal/domain/events/post"
	"LinkMe/internal/repository"
	"LinkMe/internal/repository/cache"
	"LinkMe/internal/repository/dao"
	"LinkMe/internal/service"
	"LinkMe/ioc"
	"LinkMe/utils/jwt"
)

import (
	_ "github.com/google/wire"
)

// Injectors from wire.go:

func InitWebServer() *Cmd {
	db := ioc.InitDB()
	node := ioc.InitializeSnowflakeNode()
	userDAO := dao.NewUserDAO(db, node)
	cmdable := ioc.InitRedis()
	userCache := cache.NewUserCache(cmdable)
	logger := ioc.InitLogger()
	userRepository := repository.NewUserRepository(userDAO, userCache, logger)
	userService := service.NewUserService(userRepository, logger)
	handler := jwt.NewJWTHandler(cmdable)
	userHandler := api.NewUserHandler(userService, handler, logger)
	client := ioc.InitMongoDB()
	postDAO := dao.NewPostDAO(db, logger, client)
	postCache := cache.NewPostCache(cmdable, logger)
	postRepository := repository.NewPostRepository(postDAO, logger, postCache)
	interactiveDAO := dao.NewInteractiveDAO(db, logger)
	interactiveCache := cache.NewInteractiveCache(cmdable)
	interactiveRepository := repository.NewInteractiveRepository(interactiveDAO, logger, interactiveCache)
	interactiveService := service.NewInteractiveService(interactiveRepository, logger)
	saramaClient := ioc.InitSaramaClient()
	syncProducer := ioc.InitSyncProducer(saramaClient)
	producer := post.NewSaramaSyncProducer(syncProducer)
	postService := service.NewPostService(postRepository, logger, interactiveService, producer)
	postHandler := api.NewPostHandler(postService, logger, interactiveService)
	v := ioc.InitMiddlewares(handler, logger)
	engine := ioc.InitWebServer(userHandler, postHandler, v)
	interactiveReadEventConsumer := post.NewInteractiveReadEventConsumer(interactiveRepository, saramaClient, logger)
	v2 := ioc.InitConsumers(interactiveReadEventConsumer)
	cmd := &Cmd{
		server:   engine,
		consumer: v2,
	}
	return cmd
}
