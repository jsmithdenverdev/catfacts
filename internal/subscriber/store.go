package subscriber

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type sqliteStore struct {
	db *sql.DB
}

func NewSqliteSubscriberStore(conn string) (Store, error) {
	db, err := sql.Open("sqlite3", conn)

	if err != nil {
		return nil, fmt.Errorf("could not open sqlite3 connection: %w", err)
	}

	// initialize the db
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS subscriber (
			id INTEGER PRIMARY KEY,
			contact TEXT NOT NULL
		)`,
	)

	if err != nil {
		return nil, fmt.Errorf("could not create table subscriber: %w", err)
	}

	return sqliteStore{
		db,
	}, nil
}

func (store sqliteStore) Read(contact string) (*Subscriber, error) {
	var subscriber Subscriber

	row := store.db.QueryRow("SELECT contact FROM subscriber WHERE contact = $1", contact)

	if row == nil {
		return nil, nil
	}

	err := row.Scan(&subscriber.Contact)

	// seems to be some bogus issue where the db will fail if we have an empty result set
	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("could not read subscriber from database: %w", err)
	}

	if len(subscriber.Contact) > 0 {
		return &subscriber, nil
	}

	return nil, nil
}

func (store sqliteStore) Write(subscriber Subscriber) error {
	_, err := store.db.Exec("INSERT INTO subscriber (contact) VALUES (?)", subscriber.Contact)

	if err != nil {
		return fmt.Errorf("could not Write subscriber to database: %w", err)
	}

	return nil
}

func (store sqliteStore) List() ([]Subscriber, error) {
	rows, err := store.db.Query("SELECT contact FROM subscriber")

	if err != nil {
		return nil, fmt.Errorf("could not List subscribers from database: %w", err)
	}

	defer rows.Close()

	subscribers := make([]Subscriber, 0)

	for rows.Next() {
		var subscriber Subscriber

		if err = rows.Scan(&subscriber.Contact); err != nil {
			return nil, fmt.Errorf("could not List subscribers from database: %w", err)
		}

		subscribers = append(subscribers, subscriber)
	}

	if err = rows.Close(); err != nil {
		return nil, fmt.Errorf("could not List subscribers from database: %w", err)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("could not List subscribers from database: %w", err)
	}

	return subscribers, nil
}

func (store sqliteStore) Delete(contact string) error {
	_, err := store.db.Exec("DELETE FROM subscriber WHERE contact = ?", contact)

	if err != nil {
		return fmt.Errorf("could not Delete subscriber from database: %w", err)
	}

	return nil
}
