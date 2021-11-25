package catfactninja

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type CatFactNinjaClient struct {
}

func (client CatFactNinjaClient) Retrieve() (fact string, err error) {
	url := "https://catfact.ninja/fact"
	res, err := http.Get(url)

	if err != nil {
		return
	}

	if res.Body != nil {
		defer func() {
			err = res.Body.Close()
		}()
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return
	}

	response := struct {
		Fact string `json:"fact"`
	}{}

	err = json.Unmarshal(body, &response)

	if err != nil {
		return
	}

	return response.Fact, nil
}
