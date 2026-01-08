package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Fetcher struct {
	registryURL string
	client      *http.Client
}

func NewFetcher(registryURL string) *Fetcher {
	return &Fetcher{
		registryURL: registryURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (f *Fetcher) FetchRegistry() (*Registry, error) {
	resp, err := f.client.Get(f.registryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read registry: %w", err)
	}

	var registry Registry
	if err := json.Unmarshal(body, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse registry: %w", err)
	}

	return &registry, nil
}

func (f *Fetcher) FetchFile(url string) ([]byte, error) {
	resp, err := f.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("file returned status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
