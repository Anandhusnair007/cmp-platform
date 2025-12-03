package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/cmp-platform/backend/internal/db"
)

type WebhookHandler struct {
	db *db.DB
}

func NewWebhookHandler(db *db.DB) *WebhookHandler {
	return &WebhookHandler{db: db}
}

func (h *WebhookHandler) IngestWebhook(c *gin.Context) {
	var event struct {
		EventType string                 `json:"event_type" binding:"required"`
		Timestamp string                 `json:"timestamp" binding:"required"`
		Payload   map[string]interface{} `json:"payload"`
	}

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log webhook event as audit log
	detailsJSON, _ := json.Marshal(event.Payload)
	_, err := h.db.Exec(
		`INSERT INTO audit_logs (entity_type, entity_id, action, performed_by, timestamp, details, ip_address, user_agent)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		"webhook", event.EventType, "ingest", "external", time.Now(), detailsJSON,
		c.ClientIP(), c.GetHeader("User-Agent"),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process webhook"})
		return
	}

	// TODO: Process webhook event based on event_type
	// - certificate_expiring
	// - certificate_expired
	// - agent_offline
	// - etc.

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}
