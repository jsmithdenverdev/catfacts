package sqlite3

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/jsmithdenverdev/catfacts/internal/subscriber"
)

//sqliteSubscriberStore is a subscriber store that can be used with sqlite databases
type subscriberStore struct {
	db *sql.DB
}

func (s subscriberStore) Read(contact subscriber.Contact) (subscriber.Subscriber, error) {
	// check for existing
	var sub subscriber.Subscriber
	row := s.db.QueryRow("SELECT contact FROM SUBSCRIBER WHERE contact = $1", contact)
	err := row.Scan(&sub.Contact)

	// QueryRow will return an error if scan is called and no rows are found
	// we want to continue if we ge that error, but return any others
	if err != nil && err.Error() != "sql: no rows in result set" {
		return sub, fmt.Errorf("could not check for existing subscriber: %w", err)
	}

	return sub, nil
}

func (s subscriberStore) Write(subscriber subscriber.Subscriber) error {
	// insert the contact into the subscriber table
	_, err := s.db.Exec("INSERT INTO subscriber (contact) VALUES ($1)", subscriber.Contact)

	// return an error if insertion failed
	if err != nil {
		return fmt.Errorf("could not write subscriber to database: %w", err)
	}

	return nil
}

func (s subscriberStore) List() ([]subscriber.Subscriber, error) {
	rows, err := s.db.Query("SELECT contact FROM subscriber")

	if err != nil {
		return nil, fmt.Errorf("could not list subscribers from database: %w", err)
	}

	subscribers := make([]subscriber.Subscriber, 0)

	defer rows.Close()

	for rows.Next() {
		var sub subscriber.Subscriber

		if err = rows.Scan(&sub.Contact); err != nil {
			return nil, fmt.Errorf("could not List subscribers from database: %w", err)
		}

		subscribers = append(subscribers, sub)

	}

	return subscribers, nil
}

func (s subscriberStore) Delete(contact subscriber.Contact) error {
	_, err := s.db.Exec("DELETE FROM subscriber WHERE contact = ?", contact)

	if err != nil {
		return fmt.Errorf("could not Delete subscriber from database: %w", err)
	}

	return nil
}

//NewSubscriberStore opens a connection to the given sqlite data source, initializes the database and returns a
// subscriber.Store to interact with the database.
func NewSubscriberStore(dataSourceName string) (subscriber.Store, error) {
	db, err := sql.Open("sqlite3", dataSourceName)

	if err != nil {
		return nil, fmt.Errorf("could not open sqlite3 connection: %w", err)
	}

	// initialize the db
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS subscriber (" +
		"id INTEGER PRIMARY KEY," +
		"contact TEXT NOT NULL" +
		")")

	if err != nil {
		return nil, fmt.Errorf("could not initalize subscriber table: %w", err)
	}

	return subscriberStore{
		db,
	}, nil
}
