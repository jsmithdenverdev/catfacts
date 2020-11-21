package sqlite3

import (
	"catfacts/subscription"
	"database/sql"
	"github.com/mattn/go-sqlite3"
)

type SubscriberStore struct {
	db *sql.DB
}

func NewSubscriberStore(db *sql.DB) SubscriberStore {
	return SubscriberStore{
		db,
	}
}

func (s *SubscriberStore) Insert(sub subscription.Subscriber) error {
	createStatement := `insert into subscription (contact) values ($1)`
	_, err := s.db.Exec(createStatement, sub.Contact)
	if err != nil {
		if m, ok := err.(sqlite3.Error); ok {
			if m.Code == 19 {
				return subscription.ErrSubscriptionExists
			}
		}
	}

	return err
}

func (s *SubscriberStore) Delete(contact string) error {
	deleteStatement := `delete from subscription where contact = $1`
	_, err := s.db.Exec(deleteStatement, contact)
	return err
}

func (s *SubscriberStore) All() ([]*subscription.Subscriber, error) {
	subscribers := make([]*subscription.Subscriber, 0)
	selectStatement := `select contact from subscription`
	rows, err := s.db.Query(selectStatement)

	if err != nil {
		return subscribers, err
	}

	for rows.Next() {
		var contact string
		err := rows.Scan(&contact)
		if err != nil {
			return subscribers, err
		}

		subscribers = append(subscribers, &subscription.Subscriber{
			Contact: contact,
		})
	}

	return subscribers, nil
}
