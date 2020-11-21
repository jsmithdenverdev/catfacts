package main

import (
	"catfacts/http/rest"
	"catfacts/sqlite3"
	"catfacts/subscription"
	"database/sql"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":8090", "address for incoming connections")
	conn := flag.String("conn", "catfacts.db", "db conn string")

	flag.Parse()

	db, err := sql.Open("sqlite3", *conn)
	if err != nil {
		log.Fatal(err)
	}

	subscriberStore := sqlite3.NewSubscriberStore(db)
	subscriptionService := subscription.NewService(&subscriberStore)
	handler := rest.NewHandler(subscriptionService)

	log.Fatal(http.ListenAndServe(*addr, handler))
}
