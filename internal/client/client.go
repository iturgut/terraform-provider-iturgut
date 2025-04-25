package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const HostURL string = "http://localhost:8080"

type Client struct {
	HTTPClient *http.Client
	Endpoint   string
}

func NewClient(host *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		Endpoint:   HostURL,
	}

	if host != nil {
		c.Endpoint = *host
	}

	return &c, nil
}

type Engineer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Engineer API methods ----------------------------------

func (c *Client) GetEngineers() ([]Engineer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/engineers", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", res.StatusCode)
	}

	var engineers []Engineer
	if err := json.NewDecoder(res.Body).Decode(&engineers); err != nil {
		return nil, err
	}

	return engineers, nil
}

func (c *Client) CreateEngineer(engineer Engineer) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/engineers", c.Endpoint), nil)
	if err != nil {
		return err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status code %d", res.StatusCode)
	}

	return nil
}

func (c *Client) DeleteEngineer(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/engineers/%s", c.Endpoint, id), nil)
	if err != nil {
		return err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status code %d", res.StatusCode)
	}

	return nil
}

func (c *Client) UpdateEngineer(engineer Engineer) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/engineers/%s", c.Endpoint, engineer.ID), nil)
	if err != nil {
		return err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status code %d", res.StatusCode)
	}

	return nil
}

func (c *Client) GetEngineerByID(id string) (Engineer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/engineers/id/%s", c.Endpoint, id), nil)
	if err != nil {
		return Engineer{}, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return Engineer{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Engineer{}, fmt.Errorf("API returned status code %d", res.StatusCode)
	}

	var engineer Engineer
	if err := json.NewDecoder(res.Body).Decode(&engineer); err != nil {
		return Engineer{}, err
	}

	return engineer, nil
}
