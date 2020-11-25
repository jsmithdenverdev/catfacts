package catfacts

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Fact = string

type FactRetriever = func() (Fact, error)

func DistributeFactToSubscribers(service SubscriptionService, distributor Distributor, retriever FactRetriever) error {
	errs := make([]error, 0)
	fact, err := retriever()

	subscribers, err := service.All()
	if err != nil {
		return err
	}

	for _, sub := range subscribers {
		err := distributor.Distribute(sub.Contact, fact)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		builder := strings.Builder{}
		for _, err := range errs {
			builder.WriteString(err.Error())
		}

		return errors.New(builder.String())
	}

	return nil
}

func RetrieveFactFromCatfactNinja() (fact Fact, err error) {
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

	response := struct {
		Fact string `json:"fact"`
	}{}

	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Print(err)
		return "", fmt.Errorf("could not parse body: %w", err)
	}

	return response.Fact, nil
}
