package dynamo

import (
	"catfacts"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type SubscriptionService struct {
	client *dynamodb.DynamoDB
	table  *string
}

func NewSubscriptionService(session *session.Session, table string) SubscriptionService {
	client := dynamodb.New(session)

	return SubscriptionService{
		client,
		aws.String(table),
	}
}

func (s *SubscriptionService) CreateSubscription(sub catfacts.Subscription) error {
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

func (s *SubscriptionService) DeleteSubscription(contact string) error {
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

func (s *SubscriptionService) All() ([]*catfacts.Subscription, error) {
	subs := make([]*catfacts.Subscription, 0)

	result, err := s.client.Scan(&dynamodb.ScanInput{
		TableName: s.table,
	})

	if err != nil {
		return subs, err
	}

	for _, item := range result.Items {
		sub := catfacts.Subscription{}
		err := dynamodbattribute.UnmarshalMap(item, &sub)
		if err != nil {
			return subs, err
		}
		subs = append(subs, &sub)
	}

	return subs, nil
}
