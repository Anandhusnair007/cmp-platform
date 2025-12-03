package rbac

import (
	"strings"
)

// Permission represents a permission in the RBAC system
type Permission string

const (
	// Certificate permissions
	PermissionCertRead     Permission = "cert:read"
	PermissionCertCreate   Permission = "cert:create"
	PermissionCertRevoke   Permission = "cert:revoke"
	PermissionCertRotate   Permission = "cert:rotate"
	PermissionCertDownload Permission = "cert:download"

	// Inventory permissions
	PermissionInventoryRead  Permission = "inventory:read"
	PermissionInventoryScan  Permission = "inventory:scan"
	PermissionInventoryExport Permission = "inventory:export"

	// Agent permissions
	PermissionAgentRead   Permission = "agent:read"
	PermissionAgentInstall Permission = "agent:install"
	PermissionAgentManage Permission = "agent:manage"

	// Admin permissions
	PermissionAdminRead     Permission = "admin:read"
	PermissionAdminManage   Permission = "admin:manage"
	PermissionAdminAdapters Permission = "admin:adapters"
	PermissionAdminRBAC     Permission = "admin:rbac"
	PermissionAdminUsers    Permission = "admin:users"

	// Audit permissions
	PermissionAuditRead  Permission = "audit:read"
	PermissionAuditExport Permission = "audit:export"
)

// Role represents a role with permissions
type Role struct {
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
	Inherits    []string     `json:"inherits"` // Role names to inherit from
}

// RoleRegistry manages roles and permissions
type RoleRegistry struct {
	roles map[string]*Role
}

// NewRoleRegistry creates a new role registry with default roles
func NewRoleRegistry() *RoleRegistry {
	registry := &RoleRegistry{
		roles: make(map[string]*Role),
	}

	// Define default roles
	registry.defineDefaultRoles()
	return registry
}

// defineDefaultRoles defines the default roles
func (rr *RoleRegistry) defineDefaultRoles() {
	// Admin role - full access
	rr.roles["admin"] = &Role{
		Name: "admin",
		Permissions: []Permission{
			PermissionCertRead,
			PermissionCertCreate,
			PermissionCertRevoke,
			PermissionCertRotate,
			PermissionCertDownload,
			PermissionInventoryRead,
			PermissionInventoryScan,
			PermissionInventoryExport,
			PermissionAgentRead,
			PermissionAgentInstall,
			PermissionAgentManage,
			PermissionAdminRead,
			PermissionAdminManage,
			PermissionAdminAdapters,
			PermissionAdminRBAC,
			PermissionAdminUsers,
			PermissionAuditRead,
			PermissionAuditExport,
		},
	}

	// Security role - security-focused permissions
	rr.roles["security"] = &Role{
		Name: "security",
		Permissions: []Permission{
			PermissionCertRead,
			PermissionCertRevoke,
			PermissionInventoryRead,
			PermissionInventoryScan,
			PermissionAgentRead,
			PermissionAuditRead,
			PermissionAuditExport,
		},
	}

	// Developer role - can request certificates
	rr.roles["developer"] = &Role{
		Name: "developer",
		Permissions: []Permission{
			PermissionCertRead,
			PermissionCertCreate,
			PermissionCertDownload,
			PermissionInventoryRead,
			PermissionAgentRead,
		},
	}

	// Agent role - minimal permissions for agents
	rr.roles["agent"] = &Role{
		Name: "agent",
		Permissions: []Permission{
			PermissionCertRead,
			PermissionAgentRead,
		},
	}
}

// HasPermission checks if a role has a specific permission
func (rr *RoleRegistry) HasPermission(roleNames []string, permission Permission) bool {
	checkedRoles := make(map[string]bool)

	for _, roleName := range roleNames {
		if rr.hasPermissionRecursive(roleName, permission, checkedRoles) {
			return true
		}
	}

	return false
}

// hasPermissionRecursive recursively checks permissions including inherited roles
func (rr *RoleRegistry) hasPermissionRecursive(roleName string, permission Permission, checkedRoles map[string]bool) bool {
	if checkedRoles[roleName] {
		return false // Prevent infinite loops
	}
	checkedRoles[roleName] = true

	role, exists := rr.roles[roleName]
	if !exists {
		return false
	}

	// Check direct permissions
	for _, perm := range role.Permissions {
		if perm == permission {
			return true
		}
		// Check wildcard permissions (e.g., "cert:*")
		if strings.HasSuffix(string(perm), ":*") {
			permPrefix := strings.Split(string(perm), ":")[0]
			requiredPrefix := strings.Split(string(permission), ":")[0]
			if permPrefix == requiredPrefix {
				return true
			}
		}
	}

	// Check inherited roles
	for _, inheritedRole := range role.Inherits {
		if rr.hasPermissionRecursive(inheritedRole, permission, checkedRoles) {
			return true
		}
	}

	return false
}

// GetRole returns a role by name
func (rr *RoleRegistry) GetRole(roleName string) (*Role, bool) {
	role, exists := rr.roles[roleName]
	return role, exists
}

// AddRole adds a custom role
func (rr *RoleRegistry) AddRole(role *Role) {
	rr.roles[role.Name] = role
}

// ListRoles returns all registered roles
func (rr *RoleRegistry) ListRoles() []*Role {
	roles := make([]*Role, 0, len(rr.roles))
	for _, role := range rr.roles {
		roles = append(roles, role)
	}
	return roles
}
