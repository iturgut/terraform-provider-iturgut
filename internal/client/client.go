package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	body, statusCode, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", statusCode)
	}

	var engineers []Engineer
	if err := json.Unmarshal(body, &engineers); err != nil {
		return nil, err
	}

	return engineers, nil
}

func (c *Client) CreateEngineer(engineer Engineer) error {
	body, err := json.Marshal(engineer)
	if err != nil {
		return err
	}
	
	fmt.Printf("DEBUG: Creating engineer with request body: %s\n", string(body))
	
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/engineers", c.Endpoint), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	respBody, statusCode, err := c.doRequest(req)
	if err != nil {
		return err
	}

	// Debug output
	fmt.Printf("DEBUG: API response status: %d, body: %s\n", statusCode, string(respBody))
	
	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		return fmt.Errorf("API returned status code %d", statusCode)
	}

	return nil
}

func (c *Client) DeleteEngineer(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/engineers/%s", c.Endpoint, id), nil)
	if err != nil {
		return err
	}

	respBody, statusCode, err := c.doRequest(req)
	if err != nil {
		return err
	}

	// Debug output
	fmt.Printf("DEBUG: API response status: %d, body: %s\n", statusCode, string(respBody))
	
	if statusCode != http.StatusOK {
		return fmt.Errorf("API returned status code %d", statusCode)
	}

	return nil
}

func (c *Client) UpdateEngineer(engineer Engineer) error {
	body, err := json.Marshal(engineer)
	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/engineers/%s", c.Endpoint, engineer.ID), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	respBody, statusCode, err := c.doRequest(req)
	if err != nil {
		return err
	}

	// Debug output
	fmt.Printf("DEBUG: API response status: %d, body: %s\n", statusCode, string(respBody))
	
	if statusCode != http.StatusOK {
		return fmt.Errorf("API returned status code %d", statusCode)
	}

	return nil
}

func (c *Client) GetEngineerByID(id string) (Engineer, error) {
	// First, get all engineers
	engineers, err := c.GetEngineers()
	if err != nil {
		return Engineer{}, err
	}

	// Find the engineer with the matching ID
	for _, engineer := range engineers {
		if engineer.ID == id {
			return engineer, nil
		}
	}

	return Engineer{}, fmt.Errorf("engineer with ID %s not found", id)
}

func (c *Client) doRequest(req *http.Request) ([]byte, int, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, err
	}

	return body, res.StatusCode, nil
}
