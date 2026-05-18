package application

import (
	"context"
	"log"
	"time"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
	"backend-gmao/apps/maintenance-service/internal/core/ports"
	"github.com/google/uuid"
)

// TriggerService handles automatic generation of preventive work orders
type TriggerService struct {
	otService *OrdreTravailService
	repo      ports.OrdreTravailRepository
	ticker    *time.Ticker
	quit      chan struct{}
}

func NewTriggerService(otService *OrdreTravailService, repo ports.OrdreTravailRepository) *TriggerService {
	return &TriggerService{
		otService: otService,
		repo:      repo,
		quit:      make(chan struct{}),
	}
}

// Start runs the trigger service in the background
func (s *TriggerService) Start() {
	// Check every hour (for demonstration, could be daily in production)
	s.ticker = time.NewTicker(1 * time.Hour)
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.checkAndGenerateWorkOrders()
			case <-s.quit:
				s.ticker.Stop()
				return
			}
		}
	}()
	log.Println("TriggerService started")
}

// Stop stops the trigger service
func (s *TriggerService) Stop() {
	close(s.quit)
	log.Println("TriggerService stopped")
}

func (s *TriggerService) checkAndGenerateWorkOrders() {
	// In a real scenario, this would query a "MaintenancePlan" table in Asset Service
	// to find equipment that needs preventive maintenance today.
	// For now, this is a placeholder where we would generate OTs.

	log.Println("TriggerService: checking for preventive maintenance tasks...")
	
	ctx := context.Background()
	systemUserID := uuid.Nil // Representing the system

	// Example logic: Create a dummy OT
	// (You would iterate over fetched maintenance plans)
	/*
	req := &domain.CreateOrdreTravailRequest{
		Titre:           "Preventive Maintenance - Auto Generated",
		Description:     "Scheduled maintenance check",
		TypeMaintenance: string(domain.TypePreventive),
		Priorite:        string(domain.PrioriteMoyenne),
		IDEquipement:    "some-equipement-id",
	}
	_, err := s.otService.CreateOrdreTravail(ctx, req, systemUserID)
	if err != nil {
		log.Printf("TriggerService: Failed to generate OT: %v", err)
	}
	*/
	_ = ctx
	_ = systemUserID
}
