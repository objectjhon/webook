// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"github.com/ecodeclub/webook/internal/interactive"
	baguwen "github.com/ecodeclub/webook/internal/question"
	"github.com/ecodeclub/webook/internal/question/internal/event"
	"github.com/ecodeclub/webook/internal/question/internal/repository"
	"github.com/ecodeclub/webook/internal/question/internal/repository/cache"
	"github.com/ecodeclub/webook/internal/question/internal/service"
	"github.com/ecodeclub/webook/internal/question/internal/web"
	testioc "github.com/ecodeclub/webook/internal/test/ioc"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitHandler(p event.SyncDataToSearchEventProducer, intrModule *interactive.Module) (*web.Handler, error) {
	db := testioc.InitDB()
	questionDAO := baguwen.InitQuestionDAO(db)
	ecacheCache := testioc.InitCache()
	questionCache := cache.NewQuestionECache(ecacheCache)
	repositoryRepository := repository.NewCacheRepository(questionDAO, questionCache)
	mq := testioc.InitMQ()
	interactiveEventProducer, err := event.NewInteractiveEventProducer(mq)
	if err != nil {
		return nil, err
	}
	serviceService := service.NewService(repositoryRepository, p, interactiveEventProducer)
	service2 := intrModule.Svc
	handler := web.NewHandler(serviceService, service2)
	return handler, nil
}

func InitQuestionSetHandler(p event.SyncDataToSearchEventProducer, intrModule *interactive.Module) (*web.QuestionSetHandler, error) {
	db := testioc.InitDB()
	questionSetDAO := baguwen.InitQuestionSetDAO(db)
	questionSetRepository := repository.NewQuestionSetRepository(questionSetDAO)
	mq := testioc.InitMQ()
	interactiveEventProducer, err := event.NewInteractiveEventProducer(mq)
	if err != nil {
		return nil, err
	}
	questionSetService := service.NewQuestionSetService(questionSetRepository, interactiveEventProducer, p)
	serviceService := intrModule.Svc
	questionSetHandler := web.NewQuestionSetHandler(questionSetService, serviceService)
	return questionSetHandler, nil
}

// wire.go:

var moduleSet = wire.NewSet(baguwen.InitQuestionDAO, cache.NewQuestionECache, repository.NewCacheRepository, service.NewService, web.NewHandler, baguwen.InitQuestionSetDAO, repository.NewQuestionSetRepository, service.NewQuestionSetService, web.NewQuestionSetHandler, wire.Struct(new(baguwen.Module), "*"))
