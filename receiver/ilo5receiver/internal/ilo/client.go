package ilo

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	client  *http.Client
	config  ClientConfig
	baseURL string
}

func NewClient(cfg ClientConfig) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.InsecureSkipVerify},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeout,
	}

	return &Client{
		client:  client,
		config:  cfg,
		baseURL: cfg.Endpoint,
	}
}

func (c *Client) doRequest(ctx context.Context, method, path string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.config.Username, c.config.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("OData-Version", "4.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

func (c *Client) GetSystems(ctx context.Context) ([]System, error) {
	// First get the collection to find IDs
	resp, err := c.doRequest(ctx, "GET", "/redfish/v1/Systems")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code for Systems: %d", resp.StatusCode)
	}

	var collection RedfishResponse[struct {
		ID string `json:"@odata.id"`
	}]
	if err := json.NewDecoder(resp.Body).Decode(&collection); err != nil {
		return nil, fmt.Errorf("failed to decode systems collection: %w", err)
	}

	var systems []System
	for _, member := range collection.Members {
		// Fetch individual system
		// member.ID usually contains the full URI e.g. /redfish/v1/Systems/1
		sysResp, err := c.doRequest(ctx, "GET", member.ID)
		if err != nil {
			// Log error but continue? For now, return error
			return nil, fmt.Errorf("failed to get system %s: %w", member.ID, err)
		}
		defer sysResp.Body.Close()

		var sys System
		if err := json.NewDecoder(sysResp.Body).Decode(&sys); err != nil {
			return nil, fmt.Errorf("failed to decode system: %w", err)
		}
		systems = append(systems, sys)
	}

	return systems, nil
}

func (c *Client) GetChassisIds(ctx context.Context) ([]string, error) {
	resp, err := c.doRequest(ctx, "GET", "/redfish/v1/Chassis")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code for Chassis: %d", resp.StatusCode)
	}

	var collection RedfishResponse[struct {
		ID string `json:"@odata.id"`
	}]
	if err := json.NewDecoder(resp.Body).Decode(&collection); err != nil {
		return nil, fmt.Errorf("failed to decode chassis collection: %w", err)
	}

	ids := make([]string, len(collection.Members))
	for i, m := range collection.Members {
		ids[i] = m.ID // This is the URI
	}
	return ids, nil
}

func (c *Client) GetPower(ctx context.Context, chassisURI string) (*Power, error) {
	resp, err := c.doRequest(ctx, "GET", chassisURI+"/Power")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code for Power: %d", resp.StatusCode)
	}

	var power Power
	if err := json.NewDecoder(resp.Body).Decode(&power); err != nil {
		return nil, fmt.Errorf("failed to decode power: %w", err)
	}
	return &power, nil
}

func (c *Client) GetThermal(ctx context.Context, chassisURI string) (*Thermal, error) {
	resp, err := c.doRequest(ctx, "GET", chassisURI+"/Thermal")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code for Thermal: %d", resp.StatusCode)
	}

	var thermal Thermal
	if err := json.NewDecoder(resp.Body).Decode(&thermal); err != nil {
		return nil, fmt.Errorf("failed to decode thermal: %w", err)
	}
	return &thermal, nil
}

// GetDrives fetches all drives for a system
// Returns a map of StorageID -> []Drive
func (c *Client) GetDrives(ctx context.Context, systemIDLink string) (map[string][]Drive, error) {
	// 1. Get Storage Collection
	resp, err := c.doRequest(ctx, "GET", systemIDLink+"/Storage")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code for Storage: %d", resp.StatusCode)
	}

	var collection RedfishResponse[struct {
		ID string `json:"@odata.id"`
	}]
	if err := json.NewDecoder(resp.Body).Decode(&collection); err != nil {
		return nil, fmt.Errorf("failed to decode storage collection: %w", err)
	}

	drivesMap := make(map[string][]Drive)

	// 2. Iterate over Storages
	for _, storageLink := range collection.Members {
		sResp, err := c.doRequest(ctx, "GET", storageLink.ID)
		if err != nil {
			continue // Skip failed storage
		}
		defer sResp.Body.Close()

		var storage Storage
		if err := json.NewDecoder(sResp.Body).Decode(&storage); err != nil {
			continue
		}

		var drives []Drive
		// 3. Iterate over Drives in Storage
		for _, driveLink := range storage.Drives {
			dResp, err := c.doRequest(ctx, "GET", driveLink.ID)
			if err != nil {
				continue
			}
			defer dResp.Body.Close()

			var drive Drive
			if err := json.NewDecoder(dResp.Body).Decode(&drive); err != nil {
				continue
			}
			drives = append(drives, drive)
		}
		drivesMap[storage.ID] = drives
	}

	return drivesMap, nil
}
