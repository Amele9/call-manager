package main

import (
	"log"

	"github.com/Amele9/call-manager/internal/app"
)

func main() {
	service, err := app.New()
	if err != nil {
		log.Fatalln(err)
	}

	err = service.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
