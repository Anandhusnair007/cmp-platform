package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Certificate represents a TLS certificate
type Certificate struct {
	ID            string    `json:"id" db:"id"`
	Fingerprint   string    `json:"fingerprint" db:"fingerprint"`
	Subject       string    `json:"subject" db:"subject"`
	SANs          StringArray `json:"sans" db:"sans"`
	Issuer        string    `json:"issuer" db:"issuer"`
	NotBefore     time.Time `json:"not_before" db:"not_before"`
	NotAfter      time.Time `json:"not_after" db:"not_after"`
	KeyAlgo       string    `json:"key_algo" db:"key_algo"`
	KeySize       int       `json:"key_size" db:"key_size"`
	Status        string    `json:"status" db:"status"`
	OwnerID       *string   `json:"owner_id" db:"owner_id"`
	LastScannedAt *time.Time `json:"last_scanned_at" db:"last_scanned_at"`
	Source        *string   `json:"source" db:"source"`
	CertPEM       *string   `json:"cert_pem,omitempty" db:"cert_pem"`
	ChainPEM      *string   `json:"chain_pem,omitempty" db:"chain_pem"`
	PrivateKeyRef *string   `json:"private_key_ref,omitempty" db:"private_key_ref"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Owner represents a certificate owner
type Owner struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Team      *string   `json:"team" db:"team"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Adapter represents a CA adapter configuration
type Adapter struct {
	ID        string                 `json:"id" db:"id"`
	Type      string                 `json:"type" db:"type"`
	Config    map[string]interface{} `json:"config" db:"config"`
	Enabled   bool                   `json:"enabled" db:"enabled"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" db:"updated_at"`
}

// IssuanceRequest represents a certificate issuance request
type IssuanceRequest struct {
	ID           string    `json:"id" db:"id"`
	OwnerID      string    `json:"owner_id" db:"owner_id"`
	CSR          *string   `json:"csr,omitempty" db:"csr"`
	CSRAttrs     JSONB     `json:"csr_attrs" db:"csr_attrs"`
	Status       string    `json:"status" db:"status"`
	AdapterID    string    `json:"adapter_id" db:"adapter_id"`
	RequestedAt  time.Time `json:"requested_at" db:"requested_at"`
	IssuedCertID *string   `json:"issued_cert_id" db:"issued_cert_id"`
	ErrorMessage *string   `json:"error_message,omitempty" db:"error_message"`
	CompletedAt  *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Agent represents a registered agent
type Agent struct {
	ID            string                 `json:"id" db:"id"`
	Hostname      string                 `json:"hostname" db:"hostname"`
	IP            *string                `json:"ip" db:"ip"`
	LastCheckin   time.Time              `json:"last_checkin" db:"last_checkin"`
	Capabilities  map[string]interface{} `json:"capabilities" db:"capabilities"`
	Status        string                 `json:"status" db:"status"`
	AuthTokenHash *string                `json:"-" db:"auth_token_hash"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at" db:"updated_at"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID          int64                  `json:"id" db:"id"`
	EntityType  string                 `json:"entity_type" db:"entity_type"`
	EntityID    string                 `json:"entity_id" db:"entity_id"`
	Action      string                 `json:"action" db:"action"`
	PerformedBy string                 `json:"performed_by" db:"performed_by"`
	Timestamp   time.Time              `json:"timestamp" db:"timestamp"`
	Details     map[string]interface{} `json:"details" db:"details"`
	IPAddress   *string                `json:"ip_address" db:"ip_address"`
	UserAgent   *string                `json:"user_agent" db:"user_agent"`
}

// InstallationJob represents a certificate installation job
type InstallationJob struct {
	ID           string     `json:"id" db:"id"`
	AgentID      string     `json:"agent_id" db:"agent_id"`
	CertID       string     `json:"cert_id" db:"cert_id"`
	TargetPath   string     `json:"target_path" db:"target_path"`
	ReloadCmd    *string    `json:"reload_cmd" db:"reload_cmd"`
	Status       string     `json:"status" db:"status"`
	ErrorMessage *string    `json:"error_message,omitempty" db:"error_message"`
	StartedAt    *time.Time `json:"started_at,omitempty" db:"started_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// CertRequest represents an API request for certificate issuance
type CertRequest struct {
	OwnerID       string          `json:"owner_id"`
	CommonName    string          `json:"common_name"`
	SANs          []string        `json:"sans,omitempty"`
	KeyAlgorithm  string          `json:"key_algorithm,omitempty"`
	KeySize       int             `json:"key_size,omitempty"`
	AdapterID     string          `json:"adapter_id"`
	InstallTargets []InstallTarget `json:"install_targets,omitempty"`
}

// InstallTarget represents a certificate installation target
type InstallTarget struct {
	AgentID    string  `json:"agent_id"`
	Path       string  `json:"path"`
	ReloadCmd  *string `json:"reload_cmd,omitempty"`
}

// StringArray is a custom type for PostgreSQL string arrays
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}

// JSONB is a custom type for PostgreSQL JSONB columns
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}
