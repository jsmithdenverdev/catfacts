package main

import (
	"catfacts/fact"
	"catfacts/sqlite3"
	"catfacts/subscription"
	"catfacts/twilio"
	"database/sql"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	conn := flag.String("conn", "catfacts.db", "db conn string")
	sid := flag.String("sid", "", "twilio sid")
	token := flag.String("token", "", "twilio token")
	from := flag.String("from", "", "twilio from")

	flag.Parse()

	db, err := sql.Open("sqlite3", *conn)
	if err != nil {
		log.Fatal(err)
	}

	subscriberStore := sqlite3.NewSubscriberStore(db)
	subscriptionService := subscription.NewService(&subscriberStore)
	sender := twilio.NewSmsSender(*sid, *token, *from)

	distributor := fact.NewDistributor(subscriptionService, sender)

	r := func() (string, error) {
		return "cats r neat", nil
	}

	log.Fatal(distributor.DistributeFactToSubscribers(r))
}
