package compliance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

// ImmutableAuditLog represents an immutable audit log entry
type ImmutableAuditLog struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	EntityType  string                 `json:"entity_type"`
	EntityID    string                 `json:"entity_id"`
	Action      string                 `json:"action"`
	PerformedBy string                 `json:"performed_by"`
	IPAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	Details     map[string]interface{} `json:"details"`
	Hash        string                 `json:"hash"`        // Hash of previous entry + current entry
	PreviousHash string                `json:"previous_hash"` // Hash of previous entry
}

// CreateAuditEntry creates a new immutable audit log entry
func CreateAuditEntry(entityType, entityID, action, performedBy, ipAddress, userAgent string, details map[string]interface{}, previousHash string) *ImmutableAuditLog {
	entry := &ImmutableAuditLog{
		ID:          generateID(),
		Timestamp:   time.Now().UTC(),
		EntityType:  entityType,
		EntityID:    entityID,
		Action:      action,
		PerformedBy: performedBy,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Details:     details,
		PreviousHash: previousHash,
	}

	// Calculate hash (hash of previous hash + current entry data)
	entry.Hash = calculateHash(entry, previousHash)
	return entry
}

// VerifyAuditChain verifies the integrity of audit log chain
func VerifyAuditChain(entries []*ImmutableAuditLog) bool {
	if len(entries) == 0 {
		return true
	}

	previousHash := ""
	for i, entry := range entries {
		if i > 0 {
			if entry.PreviousHash != previousHash {
				return false
			}
		}

		expectedHash := calculateHash(entry, entry.PreviousHash)
		if entry.Hash != expectedHash {
			return false
		}

		previousHash = entry.Hash
	}

	return true
}

// calculateHash calculates SHA256 hash of entry data + previous hash
func calculateHash(entry *ImmutableAuditLog, previousHash string) string {
	data := map[string]interface{}{
		"id":           entry.ID,
		"timestamp":    entry.Timestamp,
		"entity_type":  entry.EntityType,
		"entity_id":    entry.EntityID,
		"action":       entry.Action,
		"performed_by": entry.PerformedBy,
		"ip_address":   entry.IPAddress,
		"user_agent":   entry.UserAgent,
		"details":      entry.Details,
		"previous_hash": previousHash,
	}

	jsonData, _ := json.Marshal(data)
	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:])
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + hex.EncodeToString([]byte(time.Now().String())[:8])
}

// ComplianceReport generates a compliance report
type ComplianceReport struct {
	ReportID      string                 `json:"report_id"`
	GeneratedAt   time.Time              `json:"generated_at"`
	Period        CompliancePeriod       `json:"period"`
	Summary       ComplianceSummary      `json:"summary"`
	AuditLogs     []*ImmutableAuditLog   `json:"audit_logs"`
	Certificates  CertificateCompliance  `json:"certificates"`
	AccessControl AccessControlReport    `json:"access_control"`
}

type CompliancePeriod struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type ComplianceSummary struct {
	TotalAuditLogs      int `json:"total_audit_logs"`
	CertificateRequests int `json:"certificate_requests"`
	CertificateIssuances int `json:"certificate_issuances"`
	CertificateRevocations int `json:"certificate_revocations"`
	FailedOperations    int `json:"failed_operations"`
}

type CertificateCompliance struct {
	TotalCertificates     int `json:"total_certificates"`
	ExpiringWithin30Days  int `json:"expiring_within_30_days"`
	ExpiringWithin7Days   int `json:"expiring_within_7_days"`
	RevokedCertificates   int `json:"revoked_certificates"`
	NonCompliantCertificates int `json:"non_compliant_certificates"`
}

type AccessControlReport struct {
	TotalUsers        int      `json:"total_users"`
	ActiveUsers       int      `json:"active_users"`
	FailedLogins      int      `json:"failed_logins"`
	PrivilegedActions int      `json:"privileged_actions"`
	Roles             []string `json:"roles"`
}
