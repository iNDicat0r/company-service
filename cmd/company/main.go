package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iNDicat0r/company/config"
	"github.com/iNDicat0r/company/internal/app/handlers"
	"github.com/iNDicat0r/company/internal/app/middlewares"
	"github.com/iNDicat0r/company/internal/app/repositories"
	"github.com/iNDicat0r/company/internal/app/services"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	configFile := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	conf, err := config.NewConfig(*configFile)
	if err != nil {
		log.Fatalf("failed to setup config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Database.User, conf.Database.Password, conf.Database.Host, conf.Database.Port, conf.Database.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// setup repositories
	userRepo, err := repositories.NewSQLUserRepository(db)
	if err != nil {
		log.Fatalf("failed to setup user repo: %v", err)
	}

	companyRepo, err := repositories.NewSQLCompanyRepository(db)
	if err != nil {
		log.Fatalf("failed to setup company repo: %v", err)
	}

	// setup services
	userSvc, err := services.NewUserService(userRepo, conf.Global.JWTSignerKey)
	if err != nil {
		log.Fatalf("failed to setup user service: %v", err)
	}

	companySvc, err := services.NewCompanyService(companyRepo)
	if err != nil {
		log.Fatalf("failed to setup company service: %v", err)
	}

	// setup handlers
	companyHandler, err := handlers.NewCompanyHandler(companySvc)
	if err != nil {
		log.Fatalf("failed to setup company handlers: %v", err)
	}

	userHandler, err := handlers.NewUserHandler(userSvc)
	if err != nil {
		log.Fatalf("failed to setup user handlers: %v", err)
	}

	// create and group router under "/v1/"
	router := gin.Default()
	v1 := router.Group("/v1")

	// company endpoints
	v1.POST("/companies/", middlewares.AuthMiddleware(conf.Global.JWTSignerKey), companyHandler.HandleCreateCompany)
	v1.GET("/companies/:companyID", companyHandler.HandleGetCompany)
	v1.DELETE("/companies/:companyID", middlewares.AuthMiddleware(conf.Global.JWTSignerKey), companyHandler.HandleDeleteCompany)
	v1.PATCH("/companies/:companyID", middlewares.AuthMiddleware(conf.Global.JWTSignerKey), companyHandler.HandleUpdateCompany)

	// auth endpoints
	v1.POST("/auth/login", userHandler.HandleAuthenticate)
	v1.GET("/auth/introspect", middlewares.AuthMiddleware(conf.Global.JWTSignerKey), userHandler.HandleIntrospect)

	err = router.Run(fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port))
	if err != nil {
		log.Fatalf("failed to start service: %v", err)
	}
}
