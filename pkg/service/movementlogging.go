package service

import (
	"context"
	"fmt"
	"time"

	"workout-manager-service/logging"
)

type movementLoggingService struct {
	logger  logging.IshiLogger
	service MovementService
}

// NewMovementLoggingService takes an IshiLogger as a dependency and returns
// a MovementService.
func NewMovementLoggingService(logger logging.IshiLogger, s MovementService) MovementService {
	return movementLoggingService{
		logger:  logger.WithFields("service", "movement"),
		service: s,
	}
}

// Create provides informative logging when requests are made to the create
// endpoint.
func (ls movementLoggingService) Create(ctx context.Context, movementName string, categoryID string) (Movement, error) {
	defer func(begin time.Time) {
		ls.logger.Info(
			method, "Create",
			requestContext, fmt.Sprintf("%+v", ctx),
			"movementName", movementName,
			"categoryID", categoryID,
			took, time.Since(begin),
		)
	}(time.Now())
	return ls.service.Create(ctx, movementName, categoryID)
}

// Get provides informative logging when requests are made to the get
// endpoint.
func (ls movementLoggingService) Get(ctx context.Context, id string) (Movement, error) {
	defer func(begin time.Time) {
		ls.logger.Info(
			method, "Get",
			requestContext, fmt.Sprintf("%+v", ctx),
			"id", id,
			took, time.Since(begin),
		)
	}(time.Now())
	return ls.service.Get(ctx, id)
}

// List provides informative logging when requests are made to the list
// endpoint.
func (ls movementLoggingService) List(ctx context.Context, categegoryName string) ([]Movement, error) {
	defer func(begin time.Time) {
		ls.logger.Info(
			method, "List",
			requestContext, fmt.Sprintf("%+v", ctx),
			"categoryName", categegoryName,
			took, time.Since(begin),
		)
	}(time.Now())
	return ls.service.List(ctx, categegoryName)
}

// Delete provides informative logging when requests are made to the delete
// endpoint.
func (ls movementLoggingService) Delete(ctx context.Context, id string) error {
	defer func(begin time.Time) {
		ls.logger.Info(
			method, "List",
			requestContext, fmt.Sprintf("%+v", ctx),
			"id", id,
			took, time.Since(begin),
		)
	}(time.Now())
	return ls.service.Delete(ctx, id)
}
