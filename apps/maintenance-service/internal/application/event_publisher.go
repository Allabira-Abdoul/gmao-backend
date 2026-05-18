package application

import (
	"context"
	"log"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
	"backend-gmao/pkg/eventbus"
)

// MaintenanceEventPublisher wraps the EventBus to publish maintenance domain events
type MaintenanceEventPublisher struct {
	bus eventbus.EventBus
}

func NewMaintenanceEventPublisher(bus eventbus.EventBus) *MaintenanceEventPublisher {
	return &MaintenanceEventPublisher{bus: bus}
}

// PublishWorkOrderCreated publishes an event when a new work order is created.
func (p *MaintenanceEventPublisher) PublishWorkOrderCreated(ctx context.Context, ordre *domain.OrdreTravailResponse) {
	if p.bus == nil {
		return
	}
	event := eventbus.Event{
		Type:    "WorkOrderCreated",
		Payload: ordre,
	}
	if err := p.bus.Publish(ctx, "maintenance_events", "workorder.created", event); err != nil {
		log.Printf("Failed to publish WorkOrderCreated event: %v", err)
	}
}

// PublishWorkOrderUpdated publishes an event when a work order is updated.
func (p *MaintenanceEventPublisher) PublishWorkOrderUpdated(ctx context.Context, ordre *domain.OrdreTravailResponse) {
	if p.bus == nil {
		return
	}
	event := eventbus.Event{
		Type:    "WorkOrderUpdated",
		Payload: ordre,
	}
	routingKey := "workorder.updated"
	if ordre.Statut == domain.StatutTermine {
		routingKey = "workorder.completed"
		event.Type = "WorkOrderCompleted"
	}
	if err := p.bus.Publish(ctx, "maintenance_events", routingKey, event); err != nil {
		log.Printf("Failed to publish %s event: %v", event.Type, err)
	}
}
