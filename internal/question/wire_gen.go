// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package baguwen

import (
	"sync"

	"github.com/ecodeclub/ecache"
	"github.com/ecodeclub/mq-api"
	"github.com/ecodeclub/webook/internal/ai"
	"github.com/ecodeclub/webook/internal/interactive"
	"github.com/ecodeclub/webook/internal/permission"
	"github.com/ecodeclub/webook/internal/question/internal/event"
	"github.com/ecodeclub/webook/internal/question/internal/job"
	"github.com/ecodeclub/webook/internal/question/internal/repository"
	"github.com/ecodeclub/webook/internal/question/internal/repository/cache"
	"github.com/ecodeclub/webook/internal/question/internal/repository/dao"
	"github.com/ecodeclub/webook/internal/question/internal/service"
	"github.com/ecodeclub/webook/internal/question/internal/web"
	"github.com/ego-component/egorm"
	"github.com/google/wire"
	"github.com/gotomicro/ego/core/econf"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func InitModule(db *gorm.DB, intrModule *interactive.Module, ec ecache.Cache, perm *permission.Module, aiModule *ai.Module, q mq.MQ) (*Module, error) {
	questionDAO := InitQuestionDAO(db)
	questionCache := cache.NewQuestionECache(ec)
	repositoryRepository := repository.NewCacheRepository(questionDAO, questionCache)
	syncDataToSearchEventProducer, err := event.NewSyncEventProducer(q)
	if err != nil {
		return nil, err
	}
	interactiveEventProducer, err := event.NewInteractiveEventProducer(q)
	if err != nil {
		return nil, err
	}
	serviceService := service.NewService(repositoryRepository, syncDataToSearchEventProducer, interactiveEventProducer)
	questionSetDAO := InitQuestionSetDAO(db)
	questionSetRepository := repository.NewQuestionSetRepository(questionSetDAO)
	questionSetService := service.NewQuestionSetService(questionSetRepository, repositoryRepository, interactiveEventProducer, syncDataToSearchEventProducer)
	examineDAO := dao.NewGORMExamineDAO(db)
	examineRepository := repository.NewCachedExamineRepository(examineDAO)
	llmService := aiModule.Svc
	examineService := service.NewLLMExamineService(repositoryRepository, examineRepository, llmService)
	adminHandler := web.NewAdminHandler(serviceService)
	adminQuestionSetHandler := web.NewAdminQuestionSetHandler(questionSetService)
	service2 := intrModule.Svc
	service3 := perm.Svc
	handler := web.NewHandler(service2, examineService, service3, serviceService)
	questionSetHandler := web.NewQuestionSetHandler(questionSetService, examineService, service2)
	examineHandler := web.NewExamineHandler(examineService)
	knowledgeJobStarter := initKnowledgeStarter(serviceService)
	module := &Module{
		Svc:                 serviceService,
		SetSvc:              questionSetService,
		ExamSvc:             examineService,
		AdminHdl:            adminHandler,
		AdminSetHdl:         adminQuestionSetHandler,
		Hdl:                 handler,
		QsHdl:               questionSetHandler,
		ExamineHdl:          examineHandler,
		KnowledgeJobStarter: knowledgeJobStarter,
	}
	return module, nil
}

// wire.go:

var ExamineHandlerSet = wire.NewSet(web.NewExamineHandler, service.NewLLMExamineService, repository.NewCachedExamineRepository, dao.NewGORMExamineDAO)

var daoOnce = sync.Once{}

func initKnowledgeStarter(svc service.Service) *job.KnowledgeJobStarter {
	baseDir := econf.GetString("job.genKnowledge.baseDir")
	return job.NewKnowledgeJobStarter(svc, baseDir)
}

func InitTableOnce(db *gorm.DB) {
	daoOnce.Do(func() {
		err := dao.InitTables(db)
		if err != nil {
			panic(err)
		}
	})
}

func InitQuestionDAO(db *egorm.Component) dao.QuestionDAO {
	InitTableOnce(db)
	return dao.NewGORMQuestionDAO(db)
}

func InitQuestionSetDAO(db *egorm.Component) dao.QuestionSetDAO {
	InitTableOnce(db)
	return dao.NewGORMQuestionSetDAO(db)
}
