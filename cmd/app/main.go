package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"url-short/internal/app/database"
	"url-short/internal/app/routes"
	"url-short/internal/app/timer"
	"url-short/pkg/stuff"
)

func main() {
	var appClose = make(chan os.Signal)
	signal.Notify(appClose, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		stuff.RegisterEnvironment()
		database.ConnectDatabase()
		timer.Initialize()
		go routes.Register()
	}()

	<-appClose
	log.Println("Приложение закрыто.")
}
