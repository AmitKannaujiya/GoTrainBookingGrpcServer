package data

import (
	cnf "train-book/cmd/config"
	server "train-book/cmd/server"
)

func LoadData(Config *cnf.Config) {
	tarinTicketingServer := server.GetServerInstance()
	tarinTicketingServer.TrainBookingService.CreateAndLoadSeats(Config)
}
