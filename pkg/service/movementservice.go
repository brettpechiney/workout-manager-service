package service

import (
	"context"

	"workout-manager-service/logging"
)

// Movement represents a discrete movement like a panda pull, squat, or bench
// press.
type Movement struct {
	Name               string `json:"id"`
	TenantID           string `json:"tenantId"`
	MovementName       string `json:"name"`
	MovementCategoryID string `json:"movementCategoryId"`
}

// MovementService describes a service that deals with movements.
type MovementService interface {
	Create(ctx context.Context, movementName string, categoryID string) (Movement, error)
	Get(ctx context.Context, id string) (Movement, error)
	List(ctx context.Context, categoryName string) ([]Movement, error)
	Delete(ctx context.Context, id string) error
}

// NewMovementService returns a basic Service with middleware wired in.
func NewMovementService(logger logging.IshiLogger) MovementService {
	var svc MovementService
	{
		svc = NewBasicMovementService()
		svc = NewMovementLoggingService(logger, svc)
	}
	return svc
}

// NewBasicMovementService returns a naïve, stateless implementation of
// MovementService.
func NewBasicMovementService() MovementService {
	return basicMovementService{}
}

type basicMovementService struct{}

// Create adds a new Movement to the database.
func (s basicMovementService) Create(ctx context.Context, movementName string, categoryID string) (Movement, error) {
	return Movement{
		Name:               "new ID",
		TenantID:           "test tenant ID",
		MovementName:       movementName,
		MovementCategoryID: "test category ID",
	}, nil
}

// Get retrieves a Movement from the database by its UUID.
func (s basicMovementService) Get(ctx context.Context, id string) (Movement, error) {
	return Movement{
		Name:               id,
		TenantID:           "test tenant ID",
		MovementName:       "Get",
		MovementCategoryID: "test category ID",
	}, nil
}

// List retrieves all movements from the database, optionally filtering by
// category name.
// TODO: filter by category
func (s basicMovementService) List(ctx context.Context, categoryName string) ([]Movement, error) {
	return []Movement{
		{
			Name:               "1",
			TenantID:           "test tenant ID",
			MovementName:       "List",
			MovementCategoryID: "test category ID",
		},
		{
			Name:               "2",
			TenantID:           "test tenant ID",
			MovementName:       "List",
			MovementCategoryID: "test category ID",
		},
	}, nil
}

// Delete removes from the database the movement with the specified ID.
func (s basicMovementService) Delete(ctx context.Context, id string) error {
	return nil
}
