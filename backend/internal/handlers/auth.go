package handlers

import (
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/cmp-platform/backend/internal/auth"
	"github.com/cmp-platform/backend/internal/config"
	"github.com/cmp-platform/backend/internal/db"
)

type AuthHandler struct {
	db     *db.DB
	config *config.Config
}

func NewAuthHandler(db *db.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		db:     db,
		config: cfg,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	User      UserInfo  `json:"user"`
	ExpiresAt time.Time `json:"expires_at"`
}

type UserInfo struct {
	ID    string   `json:"id"`
	Email string   `json:"email"`
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
	Team  string   `json:"team"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For dev: simple password check (in production, use bcrypt)
	// Check against owners table or users table
	var userID, name, email string
	var rolesJSON []byte
	var team sql.NullString

	err := h.db.QueryRow(
		`SELECT id, name, email, COALESCE(team, ''), 
		 COALESCE('["developer"]'::jsonb, '[]'::jsonb) as roles
		 FROM owners WHERE email = $1 LIMIT 1`,
		req.Email,
	).Scan(&userID, &name, &email, &team, &rolesJSON)

	if err == sql.ErrNoRows {
		// For dev: create default user if not exists
		if req.Email == "admin@example.com" && req.Password == "admin" {
			userID = "default-owner"
			name = "Default Admin"
			email = req.Email
			rolesJSON = []byte(`["admin"]`)
			team.String = "Platform"
			team.Valid = true
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// For dev: simple password check (in production, verify hashed password)
	if req.Password != "admin" && req.Password != "password" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Parse roles
	var roles []string
	if err := json.Unmarshal(rolesJSON, &roles); err != nil {
		roles = []string{"developer"}
	}

	// Generate JWT token
	token, err := auth.GenerateToken(userID, email, roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Log audit event
	h.logAudit("user", userID, "login", userID, map[string]interface{}{
		"email": email,
		"ip":    c.ClientIP(),
	})

	c.JSON(http.StatusOK, LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		User: UserInfo{
			ID:    userID,
			Email: email,
			Name:  name,
			Roles: roles,
			Team:  team.String,
		},
	})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	cl := claims.(*auth.Claims)

	// Get user details from database
	var name, email string
	var team sql.NullString
	err := h.db.QueryRow(
		"SELECT name, email, COALESCE(team, '') FROM owners WHERE id = $1",
		cl.UserID,
	).Scan(&name, &email, &team)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, UserInfo{
		ID:    cl.UserID,
		Email: email,
		Name:  name,
		Roles: cl.Roles,
		Team:  team.String,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	cl := claims.(*auth.Claims)

	// Log audit event
	h.logAudit("user", cl.UserID, "logout", cl.UserID, map[string]interface{}{
		"ip": c.ClientIP(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// RequireAuth middleware validates JWT token
func (h *AuthHandler) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

// RequireRole middleware checks if user has required role
func (h *AuthHandler) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		cl := claims.(*auth.Claims)
		if !cl.HasRole(role) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (h *AuthHandler) logAudit(entityType, entityID, action, performedBy string, details map[string]interface{}) {
	detailsJSON, _ := json.Marshal(details)
	_, _ = h.db.Exec(
		`INSERT INTO audit_logs (entity_type, entity_id, action, performed_by, timestamp, details)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		entityType, entityID, action, performedBy, time.Now(), detailsJSON,
	)
}

// Constant time string comparison
func secureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
