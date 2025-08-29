package kolide

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
)

type Device struct {
	Name       string    `json:"name"`
	LastSeenAt time.Time `json:"last_seen_at"`
	Serial     string    `json:"serial"`
}

type DevicesResponse struct {
	Devices    []Device `json:"data"`
	Pagination struct {
		Next string `json:"next"`
	} `json:"pagination"`
}

type Client struct {
	baseUrl    string
	httpClient *httpClient
}

func NewClient(apiToken string, logger logrus.FieldLogger) *Client {
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = nil
	retryClient.RequestLogHook = func(_ retryablehttp.Logger, req *http.Request, attempt int) {
		logger.WithFields(logrus.Fields{
			"url":     req.URL.String(),
			"attempt": attempt + 1,
		}).Debug("Sending request")
	}
	retryClient.ResponseLogHook = func(_ retryablehttp.Logger, resp *http.Response) {
		if resp.StatusCode == http.StatusOK {
			return
		}

		logger.
			WithField("status", resp.Status).
			Debug("unexpected status code")
	}

	return &Client{
		baseUrl: "https://k2.kolide.com/api/v0",
		httpClient: &httpClient{
			retryClient: retryClient,
			apiToken:    apiToken,
		},
	}
}

// GetDevices retrieves all devices from the Kolide API. The returned devices is sorted by the last accessed time, with
// the most recently accessed devices first.
func (c *Client) GetDevices(ctx context.Context) ([]Device, error) {
	apiUrl, err := url.Parse(c.baseUrl + "/devices?per_page=100")
	if err != nil {
		return nil, fmt.Errorf("create Kolide API URL: %w", err)
	}

	u := apiUrl.String()
	devices := make([]Device, 0)

	for {
		resp, err := c.getPaginatedDevices(ctx, u)
		if err != nil {
			return nil, err
		}

		devices = append(devices, resp.Devices...)
		u = resp.Pagination.Next
		if u == "" {
			break
		}
	}

	slices.SortStableFunc(devices, func(a, b Device) int {
		return b.LastSeenAt.Compare(a.LastSeenAt)
	})

	return devices, nil
}

func (c *Client) getPaginatedDevices(ctx context.Context, url string) (*DevicesResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get paginated response: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	devicesResponse := &DevicesResponse{}
	if err := json.NewDecoder(resp.Body).Decode(devicesResponse); err != nil {
		return nil, fmt.Errorf("decode paginated response: %w", err)
	}

	return devicesResponse, nil
}
