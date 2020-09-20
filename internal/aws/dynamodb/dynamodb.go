package dynamodb

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jsmithdenverdev/catfacts/internal/subscriber"
)

type subscriberStore struct {
	table string
	db    *dynamodb.DynamoDB
}

func (s subscriberStore) Read(contact subscriber.Contact) (subscriber.Subscriber, error) {
	params := dynamodb.GetItemInput{
		TableName: aws.String(s.table),
		Key: map[string]*dynamodb.AttributeValue{
			"contact": {
				S: aws.String(contact),
			},
		},
	}

	result, err := s.db.GetItem(&params)

	if err != nil {
		return subscriber.Subscriber{}, err
	}

	sub := subscriber.Subscriber{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &sub)

	if err != nil {
		return subscriber.Subscriber{}, fmt.Errorf("could not unmarshal subscriber from dynamo: %w", err)
	}

	return sub, nil
}

func (s subscriberStore) Write(subscriber subscriber.Subscriber) error {
	item, err := dynamodbattribute.MarshalMap(subscriber)

	if err != nil {
		return fmt.Errorf("could not marshal subscriber into dynamodb attribute map: %w", err)
	}

	params := dynamodb.PutItemInput{
		TableName: aws.String(s.table),
		Item: item,
	}

	_, err = s.db.PutItem(&params)

	if err != nil {
		return err
	}

	return nil
}

func (s subscriberStore) List() ([]subscriber.Subscriber, error) {
	params := dynamodb.ScanInput{
		TableName: aws.String(s.table),
	}

	result, err := s.db.Scan(&params)

	if err != nil {
		return nil, err
	}

	subscribers := make([]subscriber.Subscriber, 0)

	for _, item := range result.Items {
		sub := subscriber.Subscriber{}
		err := dynamodbattribute.UnmarshalMap(item, &sub)

		if err != nil {
			return nil, fmt.Errorf("could not unmarshal subscriber from dynamo: %w", err)
		}

		subscribers = append(subscribers, sub)
	}

	return subscribers, nil
}

func (s subscriberStore) Delete(contact subscriber.Contact) error {
	params := dynamodb.DeleteItemInput{
		TableName: aws.String(s.table),
		Key: map[string]*dynamodb.AttributeValue{
			"contact": {
				S: aws.String(contact),
			},
		},
	}

	_, err := s.db.DeleteItem(&params)

	if err != nil {
		return fmt.Errorf("could not delete subscriber from dynamo: %w", err)
	}

	return nil
}

func NewSubscriberStore(table string) subscriber.Store {
	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	return subscriberStore{
		table,
		db,
	}
}
