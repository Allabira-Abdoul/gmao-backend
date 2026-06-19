package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"backend-gmao/apps/auth-service/internal/core/domain"
	"backend-gmao/pkg/discovery"
)

type userClient struct {
	registry   discovery.Registry
	httpClient *http.Client
}

// NewUserClient creates a new HTTP implementation of the UserClient.
func NewUserClient(registry discovery.Registry) *userClient {
	return &userClient{
		registry:   registry,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *userClient) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	addr, err := c.registry.Discover("user-service")
	if err != nil {
		return nil, fmt.Errorf("failed to discover user-service: %w", err)
	}

	targetURL := fmt.Sprintf("http://%s/internal/by-email?email=%s", addr, url.QueryEscape(email))

	httpReq, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("X-Internal-Service", "auth-service")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call user-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("invalid email or password")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user-service returned status: %d", resp.StatusCode)
	}

	var envelope struct {
		Status string      `json:"status"`
		Data   domain.User `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	return &envelope.Data, nil
}
