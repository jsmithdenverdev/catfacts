package ssm

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type SsmCredentialFetcher struct {
	client *ssm.SSM
}

func NewSsmCredentialFetcher(session *session.Session) SsmCredentialFetcher {
	client := ssm.New(session)

	return SsmCredentialFetcher{
		client,
	}
}

func (s SsmCredentialFetcher) FetchCredential(name string) (string, error) {
	param, err := s.client.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(true),
	})

	if err != nil {
		return "", err
	}

	return *param.Parameter.Value, nil
}
