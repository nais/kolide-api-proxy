package kolide

import (
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

type httpClient struct {
	retryClient *retryablehttp.Client
	apiToken    string
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Kolide-API-Version", "2023-05-26")

	r, err := retryablehttp.FromRequest(req)
	if err != nil {
		return nil, err
	}

	return c.retryClient.Do(r)
}
