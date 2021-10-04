package api

import (
	"fmt"
	"log"

	proxy_handler "github.com/DuckLuckBreakout/proxy/internal/pkg/proxy/handler"
	proxy_repo "github.com/DuckLuckBreakout/proxy/internal/pkg/proxy/repository"
	proxy_usecase "github.com/DuckLuckBreakout/proxy/internal/pkg/proxy/usecase"
	"github.com/DuckLuckBreakout/proxy/pkg/configer"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Start() {
	configer.InitConfig("configs/app/api_server.yaml")

	postgreSqlConn, err := sqlx.Open(
		"postgres",
		fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
			configer.AppConfig.Postgresql.User,
			configer.AppConfig.Postgresql.Password,
			configer.AppConfig.Postgresql.DBName,
			configer.AppConfig.Postgresql.Host,
			configer.AppConfig.Postgresql.Port,
			configer.AppConfig.Postgresql.Sslmode,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer postgreSqlConn.Close()
	if err := postgreSqlConn.Ping(); err != nil {
		log.Fatal(err)
	}

	proxyRepo := proxy_repo.NewRepository(postgreSqlConn)
	proxyUseCase := proxy_usecase.NewUseCase(proxyRepo)
	proxyHandler := proxy_handler.NewHandler(proxyUseCase)

	apiRouter := gin.New()
	apiRouter.Use(gin.Logger())
	apiRouter.Use(gin.Recovery())


	apiRouter.GET("/requests", proxyHandler.GetAllRequestsHandler)
	apiRouter.GET("/requests/:id", proxyHandler.GetRequestHandler)
	apiRouter.GET("/repeat/:id", proxyHandler.RepeatRequestHandler)
	apiRouter.GET("/scan/:id", proxyHandler.ScanRequestHandler)

	log.Fatal(apiRouter.Run(fmt.Sprintf(
		"%s:%s",
		configer.AppConfig.Server.Host,
		configer.AppConfig.Server.ApiPort,
	)))
}
