package dynamodb

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jsmithdenverdev/catfacts"
)

type SubscriptionStore struct {
	table  *string
	client *dynamodb.Client
	logger *log.Logger
}

type dynamoSubscription struct {
	Contact string `dynamodbav:"contact"`
}

// NewSubscriptionStore creates a new SubscriptionService using the configured
// client.
func NewSubscriptionStore(table string, logger *log.Logger) (SubscriptionStore, error) {
	conf, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(os.Getenv("AWS_REGION")))

	if err != nil {
		return SubscriptionStore{}, err
	}

	client := dynamodb.NewFromConfig(conf)

	return SubscriptionStore{
		&table,
		client,
		logger,
	}, nil
}

// Create writes the given Subscription to Dynamodb returning an error if one
// occurs.
func (store SubscriptionStore) Insert(ctx context.Context, subscription catfacts.Subscription) error {
	item, err := attributevalue.MarshalMap(dynamoSubscription{
		Contact: subscription.Contact,
	})

	if store.logger != nil {
		store.logger.Printf("subscription: %v\n", item)
	}

	if err != nil {
		return err
	}

	input := dynamodb.PutItemInput{
		TableName: store.table,
		Item:      item,
	}

	_, err = store.client.PutItem(ctx, &input)

	if err != nil {
		return err
	}

	return nil
}

// Delete removes a Subscription with the given contact from Dynamodb returning
// an error if one occurs.
func (store SubscriptionStore) Delete(ctx context.Context, contact string) error {
	input := dynamodb.DeleteItemInput{
		TableName: store.table,
		Key: map[string]types.AttributeValue{
			"contact": &types.AttributeValueMemberS{
				Value: contact,
			},
		},
	}

	_, err := store.client.DeleteItem(ctx, &input)

	if err != nil {
		return err
	}

	return nil
}

// List returns all the Subscriptions found in Dynamodb or a nil slice and
// an error if one occurs.
func (store SubscriptionStore) All(ctx context.Context) ([]catfacts.Subscription, error) {
	input := dynamodb.ScanInput{
		TableName: store.table,
	}

	results, err := store.client.Scan(ctx, &input)

	if err != nil {
		return nil, err
	}

	var subscriptions []catfacts.Subscription

	for _, item := range results.Items {
		var subscription dynamoSubscription

		if err := attributevalue.UnmarshalMap(item, &subscription); err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, catfacts.Subscription{
			Contact: subscription.Contact,
		})
	}

	return subscriptions, nil
}
