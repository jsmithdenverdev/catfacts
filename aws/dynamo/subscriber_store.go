package dynamo

import (
	"catfacts/subscription"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type SubscriberStore struct {
	client *dynamodb.DynamoDB
	table  *string
}

func NewSubscriberStore(table string) SubscriberStore {
	sess := session.Must(session.NewSession())
	client := dynamodb.New(sess)

	return SubscriberStore{
		client,
		aws.String(table),
	}
}

func (s *SubscriberStore) Insert(sub subscription.Subscriber) error {
	_, err := s.client.PutItem(&dynamodb.PutItemInput{
		TableName: s.table,
		Item: map[string]*dynamodb.AttributeValue{
			"contact": {
				S: aws.String(sub.Contact),
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriberStore) Delete(contact string) error {
	_, err := s.client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: s.table,
		Key: map[string]*dynamodb.AttributeValue{
			"contact": {
				S: aws.String(contact),
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriberStore) All() ([]*subscription.Subscriber, error) {
	subs := make([]*subscription.Subscriber, 0)

	result, err := s.client.Scan(&dynamodb.ScanInput{
		TableName: s.table,
	})

	if err != nil {
		return subs, err
	}

	for _, item := range result.Items {
		sub := subscription.Subscriber{}
		err := dynamodbattribute.UnmarshalMap(item, &sub)
		if err != nil {
			return subs, err
		}
		subs = append(subs, &sub)
	}

	return subs, nil
}
