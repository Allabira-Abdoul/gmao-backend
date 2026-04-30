package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a role in the GMAO system (e.g., Administrateur, Technicien, Manager).
type Role struct {
	IDRole      uuid.UUID       `gorm:"column:id_role;type:uuid;primaryKey;default:gen_random_uuid()" json:"id_role"`
	Libelle     string          `gorm:"column:libelle;uniqueIndex;not null" json:"libelle"`
	Description string          `gorm:"column:description" json:"description"`
	Privileges  []RolePrivilege `gorm:"foreignKey:IDRole;references:IDRole" json:"privileges,omitempty"`
	CreatedAt   time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides the default table name.
func (Role) TableName() string {
	return "roles"
}

// RolePrivilege represents the many-to-many relationship between roles and privileges.
type RolePrivilege struct {
	IDRole    uuid.UUID `gorm:"column:id_role;type:uuid;primaryKey" json:"id_role"`
	Privilege string    `gorm:"column:privilege;primaryKey" json:"privilege"`
}

// TableName overrides the default table name.
func (RolePrivilege) TableName() string {
	return "role_privileges"
}

// RoleResponse is the DTO returned by API endpoints.
type RoleResponse struct {
	IDRole      uuid.UUID `json:"id_role"`
	Libelle     string    `json:"libelle"`
	Description string    `json:"description"`
	Privileges  []string  `json:"privileges"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts a Role to a RoleResponse.
func (r *Role) ToResponse() RoleResponse {
	privileges := make([]string, 0, len(r.Privileges))
	for _, rp := range r.Privileges {
		privileges = append(privileges, rp.Privilege)
	}

	return RoleResponse{
		IDRole:      r.IDRole,
		Libelle:     r.Libelle,
		Description: r.Description,
		Privileges:  privileges,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// GetPrivilegeStrings extracts the privilege names from the role.
func (r *Role) GetPrivilegeStrings() []string {
	privileges := make([]string, 0, len(r.Privileges))
	for _, rp := range r.Privileges {
		privileges = append(privileges, rp.Privilege)
	}
	return privileges
}

// CreateRoleRequest is the DTO for creating a new role.
type CreateRoleRequest struct {
	Libelle     string   `json:"libelle" binding:"required,min=2,max=100"`
	Description string   `json:"description" binding:"max=500"`
	Privileges  []string `json:"privileges" binding:"required,min=1"`
}

// UpdateRoleRequest is the DTO for updating an existing role.
type UpdateRoleRequest struct {
	Libelle     *string `json:"libelle,omitempty" binding:"omitempty,min=2,max=100"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=500"`
}

// SetPrivilegesRequest is the DTO for setting a role's privileges.
type SetPrivilegesRequest struct {
	Privileges []string `json:"privileges" binding:"required,min=1"`
}
