package tenancy

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
)

// Tenant represents a tenant in multi-tenant system
type Tenant struct {
	ID      string                 `json:"id"`
	Name    string                 `json:"name"`
	Domain  *string                `json:"domain"`
	Enabled bool                   `json:"enabled"`
	Metadata map[string]interface{} `json:"metadata"`
}

// TenantContext stores tenant information in context
type TenantContext struct {
	TenantID string
	Tenant   *Tenant
}

type tenantKey struct{}

// WithTenant adds tenant to context
func WithTenant(ctx context.Context, tenantID string, tenant *Tenant) context.Context {
	return context.WithValue(ctx, tenantKey{}, &TenantContext{
		TenantID: tenantID,
		Tenant:   tenant,
	})
}

// GetTenantFromContext retrieves tenant from context
func GetTenantFromContext(ctx context.Context) (*TenantContext, error) {
	tc, ok := ctx.Value(tenantKey{}).(*TenantContext)
	if !ok || tc == nil {
		return nil, errors.New("tenant not found in context")
	}
	return tc, nil
}

// TenantManager manages tenants
type TenantManager struct {
	db interface {
		QueryRow(query string, args ...interface{}) *sql.Row
		Query(query string, args ...interface{}) (*sql.Rows, error)
	}
}

// NewTenantManager creates a new tenant manager
func NewTenantManager(db interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
}) *TenantManager {
	return &TenantManager{db: db}
}

// GetTenant retrieves a tenant by ID
func (tm *TenantManager) GetTenant(tenantID string) (*Tenant, error) {
	var tenant Tenant
	var domain sql.NullString
	var metadataJSON []byte

	err := tm.db.QueryRow(
		"SELECT id, name, domain, enabled, metadata FROM tenants WHERE id = $1",
		tenantID,
	).Scan(&tenant.ID, &tenant.Name, &domain, &tenant.Enabled, &metadataJSON)

	if err == sql.ErrNoRows {
		return nil, errors.New("tenant not found")
	}
	if err != nil {
		return nil, err
	}

	if domain.Valid {
		tenant.Domain = &domain.String
	}

	if len(metadataJSON) > 0 {
		json.Unmarshal(metadataJSON, &tenant.Metadata)
	}

	return &tenant, nil
}

// IsTenantEnabled checks if tenant is enabled
func (tm *TenantManager) IsTenantEnabled(tenantID string) (bool, error) {
	var enabled bool
	err := tm.db.QueryRow(
		"SELECT enabled FROM tenants WHERE id = $1",
		tenantID,
	).Scan(&enabled)

	if err == sql.ErrNoRows {
		return false, errors.New("tenant not found")
	}

	return enabled, err
}
