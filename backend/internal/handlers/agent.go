package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/cmp-platform/backend/internal/db"
	"github.com/cmp-platform/backend/internal/models"
)

type AgentHandler struct {
	db *db.DB
}

func NewAgentHandler(db *db.DB) *AgentHandler {
	return &AgentHandler{db: db}
}

func (h *AgentHandler) ListAgents(c *gin.Context) {
	rows, err := h.db.Query(
		`SELECT id, hostname, ip, last_checkin, capabilities, status, created_at, updated_at 
		 FROM agents ORDER BY last_checkin DESC`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query agents"})
		return
	}
	defer rows.Close()

	agents := []models.Agent{}
	for rows.Next() {
		var agent models.Agent
		var capabilities models.JSONB
		err := rows.Scan(
			&agent.ID, &agent.Hostname, &agent.IP, &agent.LastCheckin,
			&capabilities, &agent.Status, &agent.CreatedAt, &agent.UpdatedAt,
		)
		if err != nil {
			continue
		}
		agent.Capabilities = capabilities
		agents = append(agents, agent)
	}

	c.JSON(http.StatusOK, gin.H{
		"agents": agents,
	})
}

func (h *AgentHandler) InstallCertificate(c *gin.Context) {
	agentID := c.Param("id")

	var req struct {
		CertID    string  `json:"cert_id" binding:"required"`
		Path      string  `json:"path" binding:"required"`
		ReloadCmd *string `json:"reload_cmd"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify agent exists
	var agentExists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM agents WHERE id = $1)", agentID).Scan(&agentExists)
	if err != nil || !agentExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	// Verify certificate exists
	var certExists bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM certificates WHERE id = $1)", req.CertID).Scan(&certExists)
	if err != nil || !certExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found"})
		return
	}

	// Create installation job
	jobID := uuid.New().String()
	_, err = h.db.Exec(
		`INSERT INTO installation_jobs (id, agent_id, cert_id, target_path, reload_cmd, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		jobID, agentID, req.CertID, req.Path, req.ReloadCmd, "pending", time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create installation job"})
		return
	}

	// Log audit event
	detailsJSON, _ := json.Marshal(map[string]interface{}{
		"cert_id":    req.CertID,
		"agent_id":   agentID,
		"target_path": req.Path,
	})
	_, _ = h.db.Exec(
		`INSERT INTO audit_logs (entity_type, entity_id, action, performed_by, timestamp, details)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		"installation_job", jobID, "create", "system", time.Now(), detailsJSON,
	)

	// TODO: Notify agent via message queue or direct API call

	c.JSON(http.StatusAccepted, gin.H{
		"job_id": jobID,
		"status": "pending",
	})
}
