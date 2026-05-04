package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"backend-gmao/apps/authentication-service/internal/core/domain"
	"backend-gmao/pkg/discovery"
	"github.com/google/uuid"
)

type userClient struct {
	registry discovery.Registry
	client   *http.Client
}

func NewUserClient(registry discovery.Registry) *userClient {
	return &userClient{
		registry: registry,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *userClient) FindUserByEmail(ctx context.Context, email string) (*domain.UserInfo, error) {
	return c.callUserService(ctx, "/internal/by-email?email="+email)
}

func (c *userClient) FindUserByID(ctx context.Context, id uuid.UUID) (*domain.UserInfo, error) {
	// We need to add this endpoint to user-service internal handler if it doesn't exist.
	// Based on previous steps, I only added /internal/by-email.
	// Let me add /internal/by-id to user-service later.
	return c.callUserService(ctx, fmt.Sprintf("/internal/by-id?id=%s", id.String()))
}

func (c *userClient) callUserService(ctx context.Context, path string) (*domain.UserInfo, error) {
	// 1. Discover user-service address via Consul
	addr, err := c.registry.Discover("user-service")
	if err != nil {
		return nil, fmt.Errorf("user-service not found: %v", err)
	}

	// 2. Build target URL
	url := fmt.Sprintf("http://%s%s", addr, path)

	// 3. Prepare request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 4. Add internal service header
	req.Header.Set("X-Internal-Service", "authentication-service")

	// 5. Execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user-service returned status: %d", resp.StatusCode)
	}

	// 6. Decode response directly into a struct with a strongly typed Data field
	// ⚡ Bolt Optimization: Avoided decoding into a generic interface{}, marshaling it to JSON,
	// and then unmarshaling it again. This eliminates two unnecessary reflection passes and
	// multiple allocations per API call.
	var apiResp struct {
		Success bool             `json:"success"`
		Data    *domain.UserInfo `json:"data,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	// 7. Extract UserInfo from Data
	return apiResp.Data, nil
}
