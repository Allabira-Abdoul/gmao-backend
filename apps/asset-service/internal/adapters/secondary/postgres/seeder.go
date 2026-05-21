package postgres

import (
	"log"
	"time"

	"backend-gmao/apps/asset-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Seed populates initial asset data.
func Seed(db *gorm.DB) {
	var count int64
	db.Model(&domain.Asset{}).Count(&count)
	if count > 0 {
		return
	}

	log.Println("Seeding default asset data...")
	initialAssets := []domain.Asset{
		{
			ID:            uuid.New(),
			Name:          "Ventilation Unit A1",
			Code:          "VENT-A1",
			Status:        "OPERATIONAL",
			Category:      "HVAC",
			Location:      "Building A, Roof",
			PurchaseDate:  time.Now().AddDate(-2, 0, 0),
			PurchaseValue: 5400.0,
		},
		{
			ID:            uuid.New(),
			Name:          "Hydraulic Pump P2",
			Code:          "PUMP-P2",
			Status:        "OPERATIONAL",
			Category:      "HYDRAULIC",
			Location:      "Production Area B",
			PurchaseDate:  time.Now().AddDate(-1, -6, 0),
			PurchaseValue: 12500.0,
		},
	}

	for _, a := range initialAssets {
		db.Create(&a)
	}
	log.Println("Seeding default asset data completed")
}
