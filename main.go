package main

import (
	"log"
	conf "train-book/cmd/config"
	dataLoader "train-book/cmd/data"
	server "train-book/cmd/server"
)

func main() {
	config, err := conf.GetConfig()
	if err != nil {
		log.Fatalf("Config loading failed : %v", err)
		return
	}
	dataLoader.LoadData(config)
	server.Execute(config)
}
