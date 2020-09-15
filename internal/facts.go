package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Fact = string

type FactLoader interface {
	Load() (Fact, error)
}

type FactService struct {
	store       SubscriberStore
	loader      FactLoader
	distributor Distributor
}

type apiFactLoader struct{}

func (f FactService) SendFactToSubscribers() error {
	subscribers, err := f.store.List()

	if err != nil {
		return fmt.Errorf("could not retrieve list of subscribers: %w", err)
	}

	fact, err := f.loader.Load()

	if err != nil {
		return fmt.Errorf("could not retrieve fact: %w", err)
	}

	if len(fact) == 0 {
		return errors.New("retriever returned blank fact")
	}

	sendErrs := make([]error, 0)

	for _, sub := range subscribers {
		err := f.distributor.Distribute(sub.Contact, fmt.Sprintf("Meow! %v =^._.^=", fact))

		if err != nil {
			sendErrs = append(sendErrs, err)
		}
	}

	if len(sendErrs) > 0 {
		composite := ""

		for _, err := range sendErrs {
			composite += err.Error()
			composite += "\n"
		}

		return errors.New(composite)
	}

	return nil
}

func (a apiFactLoader) Load() (Fact, error) {
	type apiResponse struct {
		Fact string `json:"fact"`
	}

	url := "https://catfact.ninja/fact"
	res, err := http.Get(url)

	if err != nil {
		return "", fmt.Errorf("request to %s failed: %w", url, err)
	}

	if res.Body != nil {
		defer func() {
			if err := res.Body.Close(); err != nil {
				err = fmt.Errorf("could not close request body: %w", err)
			}
		}()
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Print(err)
		return "", fmt.Errorf("could not Read response body: %w", err)
	}

	response := apiResponse{}
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Print(err)
		return "", fmt.Errorf("could not parse body: %w", err)
	}

	return response.Fact, nil

}

func NewFactService(store SubscriberStore, loader FactLoader, distributor Distributor) FactService {
	return FactService{
		store,
		loader,
		distributor,
	}
}

func NewApiFactLoader() FactLoader {
	return apiFactLoader{}
}
