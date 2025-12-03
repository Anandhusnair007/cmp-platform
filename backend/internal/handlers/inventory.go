package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/cmp-platform/backend/internal/db"
	"github.com/cmp-platform/backend/internal/models"
)

type InventoryHandler struct {
	db *db.DB
}

func NewInventoryHandler(db *db.DB) *InventoryHandler {
	return &InventoryHandler{db: db}
}

func (h *InventoryHandler) GetInventory(c *gin.Context) {
	// Parse query parameters
	source := c.Query("source")
	expired := c.Query("expired") == "true"
	expiringSoonDays := 0
	if daysStr := c.Query("expiring_soon"); daysStr != "" {
		if days, err := strconv.Atoi(daysStr); err == nil {
			expiringSoonDays = days
		}
	}

	query := `SELECT id, fingerprint, subject, sans, issuer, not_before, not_after, 
	          key_algo, key_size, status, owner_id, last_scanned_at, source, 
	          created_at, updated_at 
	          FROM certificates WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if source != "" {
		query += " AND source = $" + strconv.Itoa(argIndex)
		args = append(args, source)
		argIndex++
	}

	if expired {
		query += " AND not_after < NOW()"
	}

	if expiringSoonDays > 0 {
		query += " AND not_after > NOW() AND not_after < NOW() + INTERVAL '" + strconv.Itoa(expiringSoonDays) + " days'"
	}

	query += " ORDER BY not_after ASC"

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

	c.JSON(http.StatusOK, gin.H{
		"certificates": certs,
	})
}

func (h *InventoryHandler) GetExpiringCertificates(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 30
	}

	query := `SELECT id, fingerprint, subject, sans, issuer, not_before, not_after, 
	          key_algo, key_size, status, owner_id, last_scanned_at, source, 
	          created_at, updated_at 
	          FROM certificates 
	          WHERE not_after > NOW() 
	          AND not_after < NOW() + INTERVAL '1 day' * $1
	          AND status = 'active'
	          ORDER BY not_after ASC`

	rows, err := h.db.Query(query, days)
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

	// Calculate summary
	summary := map[string]int{
		"total":       len(certs),
		"expiring_7d": 0,
		"expiring_15d": 0,
		"expiring_30d": 0,
	}

	now := time.Now()
	for _, cert := range certs {
		daysUntilExpiry := int(cert.NotAfter.Sub(now).Hours() / 24)
		if daysUntilExpiry <= 7 {
			summary["expiring_7d"]++
		}
		if daysUntilExpiry <= 15 {
			summary["expiring_15d"]++
		}
		if daysUntilExpiry <= 30 {
			summary["expiring_30d"]++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"certificates": certs,
		"summary":      summary,
	})
}
