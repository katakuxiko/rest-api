package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/katakuxiko/rest-api"
	"github.com/katakuxiko/rest-api/pkg/handlers"
	"github.com/katakuxiko/rest-api/pkg/repository"
	"github.com/katakuxiko/rest-api/pkg/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {	
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}
	if err := godotenv.Load();err != nil {
		logrus.Fatalf("error loading env variables:%s", err.Error())
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
		logrus.Fatalf("Failed to initialization db: %s",err.Error())
	}

	
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	server := new(todo.Server)

	go func() {
		if err := server.Run(viper.GetString("port"),handlers.InitRoutes()); err != nil {
			logrus.Fatalf(`error occured while running http server %s`,err.Error())
	}
	}()
	logrus.Print("App started")
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<- quit

	logrus.Print("App shutting down")

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutdown %s",err.Error())
	}
	if err := db.Close();err != nil {
		logrus.Errorf("error occured on db connection close %s",err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}