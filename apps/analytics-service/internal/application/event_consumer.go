package application

import (
	"encoding/json"
	"log"

	"backend-gmao/pkg/eventbus"
)

// EventConsumer handles incoming events from RabbitMQ and updates read models
type EventConsumer struct {
	bus eventbus.EventBus
}

func NewEventConsumer(bus eventbus.EventBus) *EventConsumer {
	return &EventConsumer{bus: bus}
}

// Start begins consuming events
func (c *EventConsumer) Start() {
	if c.bus == nil {
		log.Println("EventConsumer: No EventBus configured")
		return
	}

	err := c.bus.Subscribe("maintenance_events", "analytics_maintenance_q", "workorder.*", c.handleMaintenanceEvent)
	if err != nil {
		log.Printf("EventConsumer: Failed to subscribe to maintenance events: %v", err)
	}

	log.Println("EventConsumer started listening for events...")
}

func (c *EventConsumer) handleMaintenanceEvent(event eventbus.Event) {
	log.Printf("Analytics Service received event: %s", event.Type)
	
	switch event.Type {
	case "WorkOrderCompleted":
		c.processWorkOrderCompleted(event.Payload)
	case "WorkOrderCreated":
		c.processWorkOrderCreated(event.Payload)
	}
}

func (c *EventConsumer) processWorkOrderCompleted(payload interface{}) {
	// Parse payload into a map or DTO
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling payload: %v", err)
		return
	}

	var otData map[string]interface{}
	if err := json.Unmarshal(data, &otData); err != nil {
		log.Printf("Error unmarshalling payload to map: %v", err)
		return
	}

	log.Printf("Processing WorkOrderCompleted for OT: %v", otData["reference"])

	// 1. Update Equipment Downtime if this was a breakdown
	// 2. Add to CostAnalysis (Labor Cost = TempsPasse * HourlyRate)
	// 3. Update MTTR / MTBF KPIs in the database using gorm repo
}

func (c *EventConsumer) processWorkOrderCreated(payload interface{}) {
	log.Printf("Processing WorkOrderCreated... tracking open tickets count.")
	// Update compliance / pending orders metrics
}
