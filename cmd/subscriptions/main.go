package main

import (
	"fmt"
	"log"
	"net/http"
)

func run() error {
	a, err := createApp()

	if err != nil {
		return fmt.Errorf("could not create app: %w", err)
	}

	r := createRouter(a)

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
