package postgres

import (
	"log"
	"os"

	"backend-gmao/apps/user-service/internal/core/domain"
	"backend-gmao/pkg/auth"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Seed creates the default roles and admin user on first startup.
// This is idempotent — it skips creation if the data already exists.
func Seed(db *gorm.DB) {
	seedRoles(db)
	seedAdminUser(db)
}

func seedRoles(db *gorm.DB) {
	roles := []struct {
		Libelle     string
		Description string
		Privileges  []string
	}{
		{
			Libelle:     "Administrateur",
			Description: "Full system access — all privileges granted",
			Privileges:  domain.AllPrivileges(),
		},
		{
			Libelle:     "Manager",
			Description: "Operational management — approval, analytics, and oversight",
			Privileges: []string{
				domain.PrivilegeUserView, domain.PrivilegeUserCreate, domain.PrivilegeUserUpdate,
				domain.PrivilegeRoleView,
				domain.PrivilegeAssetView, domain.PrivilegeAssetCreate, domain.PrivilegeAssetUpdate,
				domain.PrivilegeWorkOrderView, domain.PrivilegeWorkOrderCreate, domain.PrivilegeWorkOrderUpdate,
				domain.PrivilegeWorkOrderAssign, domain.PrivilegeWorkOrderApprove, domain.PrivilegeWorkOrderClose,
				domain.PrivilegeMaintenanceView, domain.PrivilegeMaintenancePlanCreate,
				domain.PrivilegeMaintenancePlanUpdate, domain.PrivilegeMaintenanceSchedule,
				domain.PrivilegeInventoryView,
				domain.PrivilegeAnalyticsView, domain.PrivilegeAnalyticsExport,
			},
		},
		{
			Libelle:     "Technicien",
			Description: "Field technician — maintenance and asset operations",
			Privileges: []string{
				domain.PrivilegeAssetView, domain.PrivilegeAssetUpdate,
				domain.PrivilegeWorkOrderView, domain.PrivilegeWorkOrderUpdate, domain.PrivilegeWorkOrderClose,
				domain.PrivilegeMaintenanceView,
				domain.PrivilegeInventoryView, domain.PrivilegeInventoryUpdate,
			},
		},
	}

	for _, r := range roles {
		var existing domain.Role
		result := db.Where("libelle = ?", r.Libelle).First(&existing)
		if result.Error == nil {
			log.Printf("Seeder: Role '%s' already exists, skipping", r.Libelle)
			continue
		}

		role := domain.Role{
			Libelle:     r.Libelle,
			Description: r.Description,
		}

		if err := db.Create(&role).Error; err != nil {
			log.Printf("Seeder: Failed to create role '%s': %v", r.Libelle, err)
			continue
		}

		// Set privileges
		rolePrivileges := make([]domain.RolePrivilege, 0, len(r.Privileges))
		for _, p := range r.Privileges {
			rolePrivileges = append(rolePrivileges, domain.RolePrivilege{
				IDRole:    role.IDRole,
				Privilege: p,
			})
		}

		if err := db.Create(&rolePrivileges).Error; err != nil {
			log.Printf("Seeder: Failed to set privileges for role '%s': %v", r.Libelle, err)
			continue
		}

		log.Printf("Seeder: Created role '%s' with %d privileges", r.Libelle, len(r.Privileges))
	}
}

func seedAdminUser(db *gorm.DB) {
	// Check if any admin user already exists
	var count int64
	db.Model(&domain.User{}).Count(&count)
	if count > 0 {
		log.Println("Seeder: Users already exist, skipping admin user creation")
		return
	}

	// Get the Administrateur role
	var adminRole domain.Role
	if err := db.Where("libelle = ?", "Administrateur").First(&adminRole).Error; err != nil {
		log.Printf("Seeder: Cannot find Administrateur role, skipping admin user: %v", err)
		return
	}

	// Get admin password from env or use a default
	adminPassword := os.Getenv("DEFAULT_ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "Admin@2026!"
		log.Println("Seeder: WARNING — Using default admin password. Set DEFAULT_ADMIN_PASSWORD env var in production!")
	}

	hashedPassword, err := auth.HashPassword(adminPassword)
	if err != nil {
		log.Printf("Seeder: Failed to hash admin password: %v", err)
		return
	}

	adminUser := domain.User{
		IDUtilisateur: uuid.New(),
		NomComplet:    "Administrateur Système",
		Email:         "admin@gmao.local",
		MotDePasse:    hashedPassword,
		StatutCompte:  domain.StatusActive,
		IDRole:        adminRole.IDRole,
	}

	if err := db.Create(&adminUser).Error; err != nil {
		log.Printf("Seeder: Failed to create admin user: %v", err)
		return
	}

	log.Println("Seeder: Created default admin user (admin@gmao.local)")
}
