package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/cmp-platform/backend/internal/db"
	"github.com/cmp-platform/backend/internal/models"
)

type IssuerHandler struct {
	db *db.DB
}

func NewIssuerHandler(db *db.DB) *IssuerHandler {
	return &IssuerHandler{db: db}
}

func (h *IssuerHandler) RequestCertificate(c *gin.Context) {
	var req models.CertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create issuance request
	requestID := uuid.New().String()
	csrAttrs, _ := json.Marshal(req)

	query := `INSERT INTO issuance_requests (id, owner_id, csr_attrs, status, adapter_id, requested_at)
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := h.db.Exec(query, requestID, req.OwnerID, csrAttrs, "pending", req.AdapterID, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Log audit event
	h.logAudit("issuance_request", requestID, "create", "system", map[string]interface{}{
		"owner_id": req.OwnerID,
		"adapter_id": req.AdapterID,
	})

	// TODO: Trigger async processing via queue

	c.JSON(http.StatusAccepted, gin.H{
		"request_id": requestID,
		"status":     "pending",
	})
}

func (h *IssuerHandler) ListCertificates(c *gin.Context) {
	ownerID := c.Query("owner_id")
	status := c.Query("status")
	limit := c.DefaultQuery("limit", "100")
	offset := c.DefaultQuery("offset", "0")

	query := `SELECT id, fingerprint, subject, sans, issuer, not_before, not_after, 
	          key_algo, key_size, status, owner_id, last_scanned_at, source, 
	          created_at, updated_at 
	          FROM certificates WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if ownerID != "" {
		query += " AND owner_id = $" + strconv.Itoa(argIndex)
		args = append(args, ownerID)
		argIndex++
	}

	if status != "" {
		query += " AND status = $" + strconv.Itoa(argIndex)
		args = append(args, status)
		argIndex++
	}

	query += " ORDER BY created_at DESC LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query certificates"})
		return
	}
	defer rows.Close()

	certs := []models.Certificate{}
	for rows.Next() {
		var cert models.Certificate
		var sans models.StringArray
		err := rows.Scan(
			&cert.ID, &cert.Fingerprint, &cert.Subject, &sans, &cert.Issuer,
			&cert.NotBefore, &cert.NotAfter, &cert.KeyAlgo, &cert.KeySize,
			&cert.Status, &cert.OwnerID, &cert.LastScannedAt, &cert.Source,
			&cert.CreatedAt, &cert.UpdatedAt,
		)
		if err != nil {
			continue
		}
		cert.SANs = sans
		certs = append(certs, cert)
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM certificates"
	h.db.QueryRow(countQuery).Scan(&total)

	c.JSON(http.StatusOK, gin.H{
		"certificates": certs,
		"total":        total,
	})
}

func (h *IssuerHandler) GetCertificate(c *gin.Context) {
	id := c.Param("id")

	var cert models.Certificate
	var sans models.StringArray
	err := h.db.QueryRow(
		`SELECT id, fingerprint, subject, sans, issuer, not_before, not_after, 
		 key_algo, key_size, status, owner_id, last_scanned_at, source, 
		 created_at, updated_at 
		 FROM certificates WHERE id = $1`,
		id,
	).Scan(
		&cert.ID, &cert.Fingerprint, &cert.Subject, &sans, &cert.Issuer,
		&cert.NotBefore, &cert.NotAfter, &cert.KeyAlgo, &cert.KeySize,
		&cert.Status, &cert.OwnerID, &cert.LastScannedAt, &cert.Source,
		&cert.CreatedAt, &cert.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query certificate"})
		return
	}

	cert.SANs = sans
	c.JSON(http.StatusOK, cert)
}

func (h *IssuerHandler) RevokeCertificate(c *gin.Context) {
	id := c.Param("id")
	
	var reason string
	if err := c.ShouldBindJSON(&struct{ Reason *string }{&reason}); err != nil {
		reason = "unspecified"
	}

	// Update certificate status
	_, err := h.db.Exec("UPDATE certificates SET status = 'revoked', updated_at = NOW() WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke certificate"})
		return
	}

	// Log audit event
	h.logAudit("certificate", id, "revoke", "system", map[string]interface{}{
		"reason": reason,
	})

	// TODO: Call adapter to revoke certificate

	c.JSON(http.StatusOK, gin.H{
		"status": "revoked",
	})
}

func (h *IssuerHandler) logAudit(entityType, entityID, action, performedBy string, details map[string]interface{}) {
	detailsJSON, _ := json.Marshal(details)
	_, _ = h.db.Exec(
		`INSERT INTO audit_logs (entity_type, entity_id, action, performed_by, timestamp, details)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		entityType, entityID, action, performedBy, time.Now(), detailsJSON,
	)
}
