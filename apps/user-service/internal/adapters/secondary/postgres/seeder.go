package postgres

import (
	authdomain "backend-gmao/pkg/auth/domain"
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
		Name        string
		Description string
		Privileges  []string
	}{
		{
			Name:        "Administrator",
			Description: "Full system access — all privileges granted",
			Privileges:  authdomain.AllPrivileges(),
		},
		{
			Name:        "Manager",
			Description: "Operational management — approval, analytics, and oversight",
			Privileges: []string{
				authdomain.PrivilegeUserView, authdomain.PrivilegeUserCreate, authdomain.PrivilegeUserUpdate,
				authdomain.PrivilegeRoleView,
				authdomain.PrivilegeAssetView, authdomain.PrivilegeAssetCreate, authdomain.PrivilegeAssetUpdate,
				authdomain.PrivilegeWorkOrderView, authdomain.PrivilegeWorkOrderCreate, authdomain.PrivilegeWorkOrderUpdate,
				authdomain.PrivilegeWorkOrderAssign, authdomain.PrivilegeWorkOrderApprove, authdomain.PrivilegeWorkOrderClose,
				authdomain.PrivilegeMaintenanceView, authdomain.PrivilegeMaintenancePlanCreate,
				authdomain.PrivilegeMaintenancePlanUpdate, authdomain.PrivilegeMaintenanceSchedule,
				authdomain.PrivilegeInventoryView,
				authdomain.PrivilegeAnalyticsView, authdomain.PrivilegeAnalyticsExport,
			},
		},
		{
			Name:        "Technician",
			Description: "Field technician — maintenance and asset operations",
			Privileges: []string{
				authdomain.PrivilegeAssetView, authdomain.PrivilegeAssetUpdate,
				authdomain.PrivilegeWorkOrderView, authdomain.PrivilegeWorkOrderUpdate, authdomain.PrivilegeWorkOrderClose,
				authdomain.PrivilegeMaintenanceView,
				authdomain.PrivilegeInventoryView, authdomain.PrivilegeInventoryUpdate,
			},
		},
	}

	for _, r := range roles {
		var existing domain.Role
		result := db.Where("name = ?", r.Name).First(&existing)
		if result.Error == nil {
			log.Printf("Seeder: Role '%s' already exists, skipping", r.Name)
			continue
		}

		role := domain.Role{
			Name:        r.Name,
			Description: r.Description,
		}

		if err := db.Create(&role).Error; err != nil {
			log.Printf("Seeder: Failed to create role '%s': %v", r.Name, err)
			continue
		}

		// Set privileges
		rolePrivileges := make([]domain.RolePrivilege, 0, len(r.Privileges))
		for _, p := range r.Privileges {
			rolePrivileges = append(rolePrivileges, domain.RolePrivilege{
				RoleID:    role.ID,
				Privilege: p,
			})
		}

		if err := db.Create(&rolePrivileges).Error; err != nil {
			log.Printf("Seeder: Failed to set privileges for role '%s': %v", r.Name, err)
			continue
		}

		log.Printf("Seeder: Created role '%s' with %d privileges", r.Name, len(r.Privileges))
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

	// Get the Administrator role
	var adminRole domain.Role
	if err := db.Where("name = ?", "Administrator").First(&adminRole).Error; err != nil {
		log.Printf("Seeder: Cannot find Administrator role, skipping admin user: %v", err)
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
		ID:           uuid.New(),
		FullName:     "System Administrator",
		Email:        "admin@gmao.local",
		Password:     hashedPassword,
		Status:       domain.StatusActive,
		RoleID:       adminRole.ID,
	}

	if err := db.Create(&adminUser).Error; err != nil {
		log.Printf("Seeder: Failed to create admin user: %v", err)
		return
	}

	log.Println("Seeder: Created default admin user (admin@gmao.local)")
}
