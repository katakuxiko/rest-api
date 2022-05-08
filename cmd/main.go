package main

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/katakuxiko/rest-api"
	"github.com/katakuxiko/rest-api/pkg/handlers"
	"github.com/katakuxiko/rest-api/pkg/repository"
	"github.com/katakuxiko/rest-api/pkg/service"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {	
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}
	if err := godotenv.Load();err != nil {
		log.Fatalf("error loading env variables:%s", err.Error())
	}

	db,err := repository.NewPostgresDb(repository.Config{
		Host: viper.GetString("db.host"),
		Port: viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBname: viper.GetString("db.dbname"),
		SSLmode: viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatalf("Failed to initialization db: %s",err.Error())
	}

	
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	server := new(todo.Server)
	if err := server.Run(viper.GetString("8000"),handlers.InitRoutes()); err != nil {
		log.Fatalf(`error occured while running http server %s`,err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}