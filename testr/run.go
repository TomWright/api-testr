package testr

import (
	"fmt"
	"log"
	"net/http"
)

func Run(t *Test, httpClient *http.Client) error {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	log.Printf("Running test: %s\n", t.Name)

	resp, err := httpClient.Do(t.Request)
	if err != nil {
		return fmt.Errorf("could not execute request: %s", err)
	}

	for _, c := range t.Checks {
		err := c.Check(resp)
		if err != nil {
			return fmt.Errorf("failed `%T` check: %s", c, err)
		}
	}

	return nil
}
