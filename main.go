package main

import (
	"log"

	"github.com/NonsoAmadi10/p2p-analysis/app"
)

func main() {

	err := app.App().Listen("0.0.0.0:1700")

	if err != nil {
		log.Fatal(err)
	}

}
