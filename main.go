package main

import (
	"log"

	"github.com/artisomin/stickers_manager_api/actions"
)

func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
