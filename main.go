package main

import (
	"golang-microservice-template/api"
	. "golang-microservice-template/utils"
	"strconv"
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

	Log.Fatal(router.Start(":" + strconv.Itoa(port)))

	Log.Infof("[PizzaService] Started")
}
