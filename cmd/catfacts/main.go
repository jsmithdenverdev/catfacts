package main

import (
	"fmt"
	app2 "gitlab.com/jsmithdenverdev/catfacts/internal/app"
	"log"
	"net/http"
)

func run() error {
	a, err := app2.CreateApp()

	if err != nil {
		return fmt.Errorf("could not create app: %w", err)
	}

	r := CreateRouter(a)

	err = http.ListenAndServe(":8080", r)

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
