package internal

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

//Contact is a string representing a phone number.
type Contact = string

//Subscriber represents a subscriber for catfacts. A subscriber has a single field. Contact which represents the
// subscribers phone number.
type Subscriber struct {
	Contact Contact
}

type SubscriberService struct {
	store SubscriberStore
}

//SubscriberStore is an interface that provides crud operations for Subscribers.
type SubscriberStore interface {
	Read(contact Contact) (Subscriber, error)
	Write(subscriber Subscriber) error
	List() ([]Subscriber, error)
	Delete(contact Contact) error
}

//sqliteSubscriberStore is a subscriber store that can be used with sqlite databases
type sqliteSubscriberStore struct {
	db *sql.DB
}

//CreateSubscriber creates a new subscriber for a given contact. If a subscriber already
// exists with this contact an error will be returned.
func (s SubscriberService) CreateSubscriber(contact string) error {
	if len(contact) == 0 {
		return errors.New("no contact supplied")
	}

	existing, err := s.store.Read(contact)

	if err != nil {
		return fmt.Errorf("could not check for existing subscriber: %w", err)
	}

	if len(existing.Contact) > 0 {
		return fmt.Errorf("a subscriber with this contact already exists: %s", contact)
	}

	err = s.store.Write(Subscriber{Contact: contact})

	if err != nil {
		return fmt.Errorf("could not write subscriber: %w", err)
	}

	return nil
}

//DeleteSubscriber deletes a subscriber for a given contact.
func (s SubscriberService) DeleteSubscriber(contact string) error {
	err := s.store.Delete(contact)

	if err != nil {
		return fmt.Errorf("could not delete subscriber: %w", err)
	}

	return nil
}

//ListSubscribers returns a list of all subscribers.
func (s SubscriberService) ListSubscribers() ([]Subscriber, error) {
	subscribers, err := s.store.List()

	if err != nil {
		return nil, fmt.Errorf("could not list subscribers: %w", err)
	}

	return subscribers, nil
}

//Read reads a subscriber from a sqlite database using the supplied contact.
func (s sqliteSubscriberStore) Read(contact string) (Subscriber, error) {
	// check for existing
	var subscriber Subscriber
	row := s.db.QueryRow("SELECT contact FROM SUBSCRIBER WHERE contact = $1", contact)
	err := row.Scan(&subscriber.Contact)

	// QueryRow will return an error if scan is called and no rows are found
	// we want to continue if we ge that error, but return any others
	if err != nil && err.Error() != "sql: no rows in result set" {
		return subscriber, fmt.Errorf("could not check for existing subscriber: %w", err)
	}

	return subscriber, nil
}

//Write writes the given subscriber to the database. In the case of sqlite it will write the subscribers contact to the
// subscriber table.
func (s sqliteSubscriberStore) Write(subscriber Subscriber) error {
	// insert the contact into the subscriber table
	_, err := s.db.Exec("INSERT INTO subscriber (contact) VALUES ($1)", subscriber.Contact)

	// return an error if insertion failed
	if err != nil {
		return fmt.Errorf("could not write subscriber to database: %w", err)
	}

	return nil
}

//List lists the subscribers in the database.
func (s sqliteSubscriberStore) List() ([]Subscriber, error) {
	rows, err := s.db.Query("SELECT contact FROM subscriber")

	if err != nil {
		return nil, fmt.Errorf("could not list subscribers from database: %w", err)
	}

	subscribers := make([]Subscriber, 0)

	defer rows.Close()

	for rows.Next() {
		var subscriber Subscriber

		if err = rows.Scan(&subscriber.Contact); err != nil {
			return nil, fmt.Errorf("could not List subscribers from database: %w", err)
		}

		subscribers = append(subscribers, subscriber)

	}

	return subscribers, nil
}

//Delete deletes a subscriber from the database that matches the given contact.
func (s sqliteSubscriberStore) Delete(contact Contact) error {
	_, err := s.db.Exec("DELETE FROM subscriber WHERE contact = ?", contact)

	if err != nil {
		return fmt.Errorf("could not Delete subscriber from database: %w", err)
	}

	return nil
}

//NewSqliteSubscriberStore opens a connection to the given sqlite data source, initializes the database and returns a
// SubscriberStore to interact with the database.
func NewSqliteSubscriberStore(dataSourceName string) (SubscriberStore, error) {
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

	return sqliteSubscriberStore{
		db,
	}, nil
}

//NewSubscriberService creates a new SubscriberService with the given SubscriberStore
func NewSubscriberService(store SubscriberStore) SubscriberService {
	return SubscriberService{store}
}
