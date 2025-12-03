package automation

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/cmp-platform/backend/internal/db"
)

// RenewalScheduler manages automated certificate renewal
type RenewalScheduler struct {
	db           *db.DB
	checkInterval time.Duration
	renewalWindow time.Duration // Days before expiry to renew
}

// NewRenewalScheduler creates a new renewal scheduler
func NewRenewalScheduler(db *db.DB, checkInterval, renewalWindow time.Duration) *RenewalScheduler {
	return &RenewalScheduler{
		db:           db,
		checkInterval: checkInterval,
		renewalWindow: renewalWindow,
	}
}

// Start starts the renewal scheduler
func (rs *RenewalScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(rs.checkInterval)
	defer ticker.Stop()

	// Run initial check
	rs.checkAndRenew()

	for {
		select {
		case <-ctx.Done():
			log.Println("Renewal scheduler stopped")
			return
		case <-ticker.C:
			rs.checkAndRenew()
		}
	}
}

// checkAndRenew checks for expiring certificates and initiates renewal
func (rs *RenewalScheduler) checkAndRenew() {
	renewalDate := time.Now().Add(rs.renewalWindow)

	query := `SELECT id, subject, sans, issuer, owner_id, adapter_id 
	          FROM certificates 
	          WHERE status = 'active' 
	          AND not_after > NOW() 
	          AND not_after <= $1 
	          AND id NOT IN (SELECT DISTINCT cert_id FROM issuance_requests WHERE status = 'pending' OR status = 'processing')
	          ORDER BY not_after ASC`

	rows, err := rs.db.Query(query, renewalDate)
	if err != nil {
		log.Printf("Error querying expiring certificates: %v", err)
		return
	}
	defer rows.Close()

	type CertInfo struct {
		ID        string
		Subject   string
		SANs      string
		Issuer    string
		OwnerID   string
		AdapterID string
	}

	certsToRenew := []CertInfo{}
	for rows.Next() {
		var cert CertInfo
		err := rows.Scan(&cert.ID, &cert.Subject, &cert.SANs, &cert.Issuer, &cert.OwnerID, &cert.AdapterID)
		if err != nil {
			continue
		}
		certsToRenew = append(certsToRenew, cert)
	}

	for _, cert := range certsToRenew {
		log.Printf("Scheduling renewal for certificate: %s (%s)", cert.ID, cert.Subject)
		rs.scheduleRenewal(cert)
	}

	if len(certsToRenew) > 0 {
		log.Printf("Scheduled %d certificate renewals", len(certsToRenew))
	}
}

// scheduleRenewal schedules a certificate renewal
func (rs *RenewalScheduler) scheduleRenewal(cert CertInfo) {
	// Create renewal request in database
	query := `INSERT INTO issuance_requests (id, owner_id, csr_attrs, status, adapter_id, requested_at)
	          VALUES ($1, $2, $3, $4, $5, $6)`

	csrAttrs := map[string]interface{}{
		"common_name": cert.Subject,
		"renewal":     true,
		"original_cert_id": cert.ID,
	}

	csrAttrsJSON, _ := json.Marshal(csrAttrs)
	requestID := generateRequestID()

	_, err := rs.db.Exec(query, requestID, cert.OwnerID, csrAttrsJSON, "pending", cert.AdapterID, time.Now())
	if err != nil {
		log.Printf("Error creating renewal request: %v", err)
		return
	}

	// Log audit event
	auditQuery := `INSERT INTO audit_logs (entity_type, entity_id, action, performed_by, timestamp, details)
	               VALUES ($1, $2, $3, $4, $5, $6)`
	
	details := map[string]interface{}{
		"request_id": requestID,
		"cert_id":    cert.ID,
		"reason":     "automated_renewal",
	}
	detailsJSON, _ := json.Marshal(details)

	_, _ = rs.db.Exec(auditQuery, "certificate", cert.ID, "renewal_scheduled", "system", time.Now(), detailsJSON)

	// TODO: Trigger async renewal process via message queue
	log.Printf("Renewal request created: %s for cert %s", requestID, cert.ID)
}

func generateRequestID() string {
	return "renewal-" + time.Now().Format("20060102150405") + "-" + uuid.New().String()[:8]
}
