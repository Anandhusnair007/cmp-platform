package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ServiceNowClient handles ServiceNow integration
type ServiceNowClient struct {
	baseURL      string
	username     string
	password     string
	client       *http.Client
}

// ServiceNowConfig represents ServiceNow configuration
type ServiceNowConfig struct {
	BaseURL  string
	Username string
	Password string
}

// NewServiceNowClient creates a new ServiceNow client
func NewServiceNowClient(config ServiceNowConfig) *ServiceNowClient {
	return &ServiceNowClient{
		baseURL:  config.BaseURL,
		username: config.Username,
		password: config.Password,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Ticket represents a ServiceNow ticket
type Ticket struct {
	Number      string `json:"number"`
	ShortDescription string `json:"short_description"`
	Description string `json:"description"`
	State       string `json:"state"`
	Priority    string `json:"priority"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
}

// CreateTicket creates a ServiceNow ticket
func (c *ServiceNowClient) CreateTicket(ticket Ticket) (*Ticket, error) {
	url := fmt.Sprintf("%s/api/now/table/incident", c.baseURL)

	payload := map[string]interface{}{
		"short_description": ticket.ShortDescription,
		"description":       ticket.Description,
		"priority":          ticket.Priority,
		"category":          ticket.Category,
		"subcategory":       ticket.Subcategory,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ServiceNow API error: %d", resp.StatusCode)
	}

	var result struct {
		Result Ticket `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Result, nil
}

// UpdateTicket updates a ServiceNow ticket
func (c *ServiceNowClient) UpdateTicket(ticketNumber string, updates map[string]interface{}) error {
	url := fmt.Sprintf("%s/api/now/table/incident/%s", c.baseURL, ticketNumber)

	jsonData, err := json.Marshal(updates)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ServiceNow API error: %d", resp.StatusCode)
	}

	return nil
}

// CreateCertificateRequestTicket creates a ticket for certificate request
func (c *ServiceNowClient) CreateCertificateRequestTicket(requestID, commonName, requesterEmail string) (*Ticket, error) {
	return c.CreateTicket(Ticket{
		ShortDescription: fmt.Sprintf("Certificate Request: %s", commonName),
		Description: fmt.Sprintf(
			"Certificate request submitted through CMP Platform.\n\n"+
				"Request ID: %s\n"+
				"Common Name: %s\n"+
				"Requester: %s\n\n"+
				"Please review and approve this certificate request.",
			requestID, commonName, requesterEmail,
		),
		Priority: "3",
		Category: "Security",
		Subcategory: "Certificate Management",
	})
}
