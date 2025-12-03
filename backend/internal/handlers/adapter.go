package handlers

import (
	"encoding/json"
	"log"

	"github.com/cmp-platform/backend/internal/config"
	"github.com/cmp-platform/backend/internal/db"
)

type AdapterHandler struct {
	db     *db.DB
	config *config.Config
}

func NewAdapterHandler(db *db.DB, cfg *config.Config) *AdapterHandler {
	return &AdapterHandler{
		db:     db,
		config: cfg,
	}
}

// ProcessIssuanceRequest processes a pending issuance request
func (h *AdapterHandler) ProcessIssuanceRequest(requestID string) error {
	// Get request from database
	var adapterID string
	var csrAttrsJSON []byte
	err := h.db.QueryRow(
		"SELECT adapter_id, csr_attrs FROM issuance_requests WHERE id = $1",
		requestID,
	).Scan(&adapterID, &csrAttrsJSON)
	if err != nil {
		return err
	}

	// Get adapter config
	var adapterType string
	var adapterConfigJSON []byte
	err = h.db.QueryRow(
		"SELECT type, config FROM adapters WHERE id = $1 AND enabled = true",
		adapterID,
	).Scan(&adapterType, &adapterConfigJSON)
	if err != nil {
		return err
	}

	var adapterConfig map[string]interface{}
	if err := json.Unmarshal(adapterConfigJSON, &adapterConfig); err != nil {
		return err
	}

	// Parse CSR attributes
	var csrAttrs map[string]interface{}
	if err := json.Unmarshal(csrAttrsJSON, &csrAttrs); err != nil {
		return err
	}

	log.Printf("Processing issuance request %s with adapter %s (%s)", requestID, adapterID, adapterType)

	// TODO: Call appropriate adapter based on type
	// For now, just log and mark as processing
	_, err = h.db.Exec(
		"UPDATE issuance_requests SET status = $1, updated_at = NOW() WHERE id = $2",
		"processing", requestID,
	)
	return err
}
