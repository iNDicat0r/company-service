package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/iNDicat0r/company/config"
	"github.com/iNDicat0r/company/internal/app/models"
	"github.com/iNDicat0r/company/internal/app/utils"
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
	db, err := gorm.Open(mysql.Open(dsn), nil)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Company{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	hashPass, _ := utils.HashPassword("124")

	// user 1
	db.Create(&models.User{
		Name:     "Mobin",
		Username: "iNDicat0r",
		Password: hashPass,
	})
}
