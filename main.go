package main

import (
	"log"

	"github.com/NonsoAmadi10/p2p-analysis/app"
	"github.com/NonsoAmadi10/p2p-analysis/services"
)

func main() {
	go services.ConnectionMetrics()

	err := app.App().Listen("0.0.0.0:1700")

	if err != nil {
		log.Fatal(err)
	}

}
