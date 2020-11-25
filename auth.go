package catfacts

import (
	"encoding/base64"
)

type BasicCredentials struct {
	Username string
	Password string
}

type CredentialFetcher interface {
	FetchCredential(name string) (string, error)
}

func AuthorizeBasicHeader(fetcher CredentialFetcher, header, name string) error {
	// header is base64 encoded and must be decoded
	decoded, err := base64.StdEncoding.DecodeString(header)
	if err != nil {
		return err
	}

	owned, err := fetcher.FetchCredential(name)
	if err != nil {
		return err
	}

	if owned != string(decoded) {
		return ErrUnauthorized
	}

	return nil
}
