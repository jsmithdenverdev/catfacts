package fact

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func RetrieveFactFromApi() (Fact, error) {
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
