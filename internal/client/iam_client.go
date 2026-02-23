package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// IAMClientConfig configures the IAM platform client.
type IAMClientConfig struct {
	BaseURL string
	Timeout time.Duration
}

// IAMClient is an HTTP client for the IAM platform service.
type IAMClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewIAMClient creates a new IAM client.
func NewIAMClient(cfg IAMClientConfig) *IAMClient {
	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}
	return &IAMClient{
		baseURL:    cfg.BaseURL,
		httpClient: &http.Client{Timeout: cfg.Timeout},
	}
}

// UserRoleDTO represents a user role from the IAM platform.
type UserRoleDTO struct {
	ID             string  `json:"id"`
	UserID         string  `json:"user_id"`
	RoleID         string  `json:"role_id"`
	RoleName       string  `json:"role_name"`
	SchoolID       *string `json:"school_id,omitempty"`
	AcademicUnitID *string `json:"academic_unit_id,omitempty"`
}

// GrantRoleRequest represents a request to grant a role.
type GrantRoleRequest struct {
	RoleID         string  `json:"role_id"`
	SchoolID       *string `json:"school_id,omitempty"`
	AcademicUnitID *string `json:"academic_unit_id,omitempty"`
}

// GetUserRoles retrieves roles for a user from the IAM platform.
func (c *IAMClient) GetUserRoles(ctx context.Context, token, userID string) ([]UserRoleDTO, error) {
	url := fmt.Sprintf("%s/v1/users/%s/roles", c.baseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling IAM service: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IAM service error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var roles []UserRoleDTO
	if err := json.Unmarshal(body, &roles); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return roles, nil
}

// GrantRole grants a role to a user via the IAM platform.
func (c *IAMClient) GrantRole(ctx context.Context, token, userID string, grantReq GrantRoleRequest) (*UserRoleDTO, error) {
	url := fmt.Sprintf("%s/v1/users/%s/roles", c.baseURL, userID)

	bodyBytes, err := json.Marshal(grantReq)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling IAM service: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("IAM service error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var role UserRoleDTO
	if err := json.Unmarshal(body, &role); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &role, nil
}

// RevokeRole revokes a role from a user via the IAM platform.
func (c *IAMClient) RevokeRole(ctx context.Context, token, userID, roleID string) error {
	url := fmt.Sprintf("%s/v1/users/%s/roles/%s", c.baseURL, userID, roleID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("calling IAM service: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("IAM service error: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
