package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

// OIDCConfig represents OIDC provider configuration
type OIDCConfig struct {
	IssuerURL      string
	ClientID       string
	ClientSecret   string
	RedirectURL    string
	Scopes         []string
	Provider       *oauth2.Config
}

// NewOIDCConfig creates a new OIDC configuration
func NewOIDCConfig(issuerURL, clientID, clientSecret, redirectURL string) (*OIDCConfig, error) {
	// Discover OIDC configuration
	wellKnownURL := fmt.Sprintf("%s/.well-known/openid-configuration", issuerURL)
	resp, err := http.Get(wellKnownURL)
	if err != nil {
		return nil, fmt.Errorf("failed to discover OIDC config: %w", err)
	}
	defer resp.Body.Close()

	var discovery struct {
		AuthorizationEndpoint string `json:"authorization_endpoint"`
		TokenEndpoint         string `json:"token_endpoint"`
		UserInfoEndpoint      string `json:"userinfo_endpoint"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&discovery); err != nil {
		return nil, fmt.Errorf("failed to parse OIDC discovery: %w", err)
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  discovery.AuthorizationEndpoint,
			TokenURL: discovery.TokenEndpoint,
		},
	}

	return &OIDCConfig{
		IssuerURL:   issuerURL,
		ClientID:    clientID,
		ClientSecret: clientSecret,
		RedirectURL: redirectURL,
		Scopes:      []string{"openid", "profile", "email"},
		Provider:    config,
	}, nil
}

// GenerateState generates a random state for OAuth flow
func GenerateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ExchangeCode exchanges authorization code for token
func (c *OIDCConfig) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.Provider.Exchange(ctx, code)
}

// GetUserInfo retrieves user information from OIDC provider
func (c *OIDCConfig) GetUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error) {
	client := c.Provider.Client(ctx, token)
	
	// Get userinfo endpoint from discovery
	wellKnownURL := fmt.Sprintf("%s/.well-known/openid-configuration", c.IssuerURL)
	resp, err := http.Get(wellKnownURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var discovery struct {
		UserInfoEndpoint string `json:"userinfo_endpoint"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&discovery); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", discovery.UserInfoEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %s", string(body))
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// UserInfo represents OIDC user information
type UserInfo struct {
	Sub           string   `json:"sub"`
	Email         string   `json:"email"`
	EmailVerified bool     `json:"email_verified"`
	Name          string   `json:"name"`
	GivenName     string   `json:"given_name"`
	FamilyName    string   `json:"family_name"`
	Groups        []string `json:"groups"`
	Roles         []string `json:"roles"`
	Picture       string   `json:"picture"`
}
