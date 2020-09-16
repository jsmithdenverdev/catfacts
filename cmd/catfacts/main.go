package main

import (
	"fmt"
	"github.com/jsmithdenverdev/catfacts/cmd/catfacts/app"
	"log"
	"net/http"
)

func run() error {
	a, err := app.CreateApp()

	if err != nil {
		return fmt.Errorf("could not create app: %w", err)
	}

	r := app.CreateRouter(a)

	err = http.ListenAndServe(":5520", r)

	if err != nil {
		return fmt.Errorf("could not start http server: %w", err)
	}

	return nil
}

func main() {
	err := run()

	if err != nil {
		log.Fatal(err)
	}
}
