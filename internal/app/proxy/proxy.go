package proxy

import (
	"fmt"
	"log"
	"net/http"

	proxy_handler "github.com/DuckLuckBreakout/proxy/internal/pkg/proxy/handler"
	proxy_repo "github.com/DuckLuckBreakout/proxy/internal/pkg/proxy/repository"
	proxy_usecase "github.com/DuckLuckBreakout/proxy/internal/pkg/proxy/usecase"
	"github.com/DuckLuckBreakout/proxy/pkg/configer"

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

	server := &http.Server{
		Addr: fmt.Sprintf(
			"%s:%s",
			configer.AppConfig.Server.Host,
			configer.AppConfig.Server.Port,
		),
		Handler: http.HandlerFunc(proxyHandler.HandleRequestHttp),
	}

	fmt.Println("Proxy ready")
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
