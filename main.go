package main

import (
	"golang-microservice-template/api"
	. "golang-microservice-template/utils"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const (
	defaultPort = 8080
	localPort   = 8081
)

func main() {
	Log.Infof("[PizzaService] Start")

	router := api.NewRouter()

	port := defaultPort
	if Environment() == ENV_LOCAL {
		port = localPort
	}

	go func() {
		Log.Fatal(router.Start(":" + strconv.Itoa(port)))
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT)
	<-done

	router.Shutdown()
}
