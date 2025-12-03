package workflows

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/cmp-platform/backend/internal/db"
)

// ApprovalWorkflow represents an approval workflow
type ApprovalWorkflow struct {
	ID            string                 `json:"id"`
	RequestID     string                 `json:"request_id"`
	EntityType    string                 `json:"entity_type"` // "certificate", "revocation", etc.
	EntityID      string                 `json:"entity_id"`
	Status        string                 `json:"status"` // "pending", "approved", "rejected"
	RequesterID   string                 `json:"requester_id"`
	Approvers     []Approver             `json:"approvers"`
	RequiredApprovals int                `json:"required_approvals"`
	CurrentStep   int                    `json:"current_step"`
	CreatedAt     time.Time              `json:"created_at"`
	CompletedAt   *time.Time             `json:"completed_at"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// Approver represents an approver in the workflow
type Approver struct {
	UserID      string     `json:"user_id"`
	Role        string     `json:"role"`
	Status      string     `json:"status"` // "pending", "approved", "rejected"
	ApprovedAt  *time.Time `json:"approved_at"`
	RejectedAt  *time.Time `json:"rejected_at"`
	Comments    string     `json:"comments"`
	Step        int        `json:"step"`
}

// ApprovalWorkflowManager manages approval workflows
type ApprovalWorkflowManager struct {
	db *db.DB
}

// NewApprovalWorkflowManager creates a new workflow manager
func NewApprovalWorkflowManager(db *db.DB) *ApprovalWorkflowManager {
	return &ApprovalWorkflowManager{db: db}
}

// CreateWorkflow creates a new approval workflow
func (awm *ApprovalWorkflowManager) CreateWorkflow(entityType, entityID, requesterID string, approvers []Approver, requiredApprovals int, metadata map[string]interface{}) (*ApprovalWorkflow, error) {
	workflowID := uuid.New().String()
	requestID := uuid.New().String()

	approversJSON, _ := json.Marshal(approvers)
	metadataJSON, _ := json.Marshal(metadata)

	query := `INSERT INTO approval_workflows 
	          (id, request_id, entity_type, entity_id, status, requester_id, 
	           approvers, required_approvals, current_step, metadata, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := awm.db.Exec(query,
		workflowID, requestID, entityType, entityID, "pending",
		requesterID, approversJSON, requiredApprovals, 0, metadataJSON, time.Now(),
	)

	if err != nil {
		return nil, err
	}

	return &ApprovalWorkflow{
		ID:               workflowID,
		RequestID:        requestID,
		EntityType:       entityType,
		EntityID:         entityID,
		Status:           "pending",
		RequesterID:      requesterID,
		Approvers:        approvers,
		RequiredApprovals: requiredApprovals,
		CurrentStep:      0,
		CreatedAt:        time.Now(),
		Metadata:         metadata,
	}, nil
}

// Approve approves a workflow step
func (awm *ApprovalWorkflowManager) Approve(workflowID, approverID string, comments string) error {
	workflow, err := awm.GetWorkflow(workflowID)
	if err != nil {
		return err
	}

	// Find approver
	for i, approver := range workflow.Approvers {
		if approver.UserID == approverID && approver.Status == "pending" {
			workflow.Approvers[i].Status = "approved"
			now := time.Now()
			workflow.Approvers[i].ApprovedAt = &now
			workflow.Approvers[i].Comments = comments

			// Check if we have enough approvals
			approvedCount := 0
			for _, app := range workflow.Approvers {
				if app.Status == "approved" {
					approvedCount++
				}
			}

			if approvedCount >= workflow.RequiredApprovals {
				workflow.Status = "approved"
				now := time.Now()
				workflow.CompletedAt = &now
			}

			// Update database
			approversJSON, _ := json.Marshal(workflow.Approvers)
			_, err = awm.db.Exec(
				`UPDATE approval_workflows 
				 SET approvers = $1, status = $2, completed_at = $3 
				 WHERE id = $4`,
				approversJSON, workflow.Status, workflow.CompletedAt, workflowID,
			)

			return err
		}
	}

	return fmt.Errorf("approver not found or already processed")
}

// Reject rejects a workflow
func (awm *ApprovalWorkflowManager) Reject(workflowID, approverID string, comments string) error {
	workflow, err := awm.GetWorkflow(workflowID)
	if err != nil {
		return err
	}

	for i, approver := range workflow.Approvers {
		if approver.UserID == approverID {
			workflow.Approvers[i].Status = "rejected"
			now := time.Now()
			workflow.Approvers[i].RejectedAt = &now
			workflow.Approvers[i].Comments = comments
			workflow.Status = "rejected"
			now := time.Now()
			workflow.CompletedAt = &now

			approversJSON, _ := json.Marshal(workflow.Approvers)
			_, err = awm.db.Exec(
				`UPDATE approval_workflows 
				 SET approvers = $1, status = $2, completed_at = $3 
				 WHERE id = $4`,
				approversJSON, workflow.Status, workflow.CompletedAt, workflowID,
			)

			return err
		}
	}

	return fmt.Errorf("approver not found")
}

// GetWorkflow retrieves a workflow by ID
func (awm *ApprovalWorkflowManager) GetWorkflow(workflowID string) (*ApprovalWorkflow, error) {
	var workflow ApprovalWorkflow
	var approversJSON, metadataJSON []byte

	err := awm.db.QueryRow(
		`SELECT id, request_id, entity_type, entity_id, status, requester_id,
		 approvers, required_approvals, current_step, metadata, created_at, completed_at
		 FROM approval_workflows WHERE id = $1`,
		workflowID,
	).Scan(
		&workflow.ID, &workflow.RequestID, &workflow.EntityType, &workflow.EntityID,
		&workflow.Status, &workflow.RequesterID, &approversJSON,
		&workflow.RequiredApprovals, &workflow.CurrentStep, &metadataJSON,
		&workflow.CreatedAt, &workflow.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("workflow not found")
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(approversJSON, &workflow.Approvers)
	json.Unmarshal(metadataJSON, &workflow.Metadata)

	return &workflow, nil
}
