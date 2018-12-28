package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"workout-manager-service/pkg/service"
)

// MovementSet is a helper struct that collects all of the Movement endpoints
// in the workout manager service.
type MovementSet struct {
	CreateEndpoint endpoint.Endpoint
	GetEndpoint    endpoint.Endpoint
	ListEndpoint   endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
}

// NewMovementSet returns a MovementSet that wraps the provided
// MovementService and wires in the endpoint middleware.
func NewMovementSet(svc service.MovementService) MovementSet {
	return MovementSet{
		CreateEndpoint: MakeCreateMovementEndpoint(svc),
		GetEndpoint:    MakeGetMovementEndpoint(svc),
		ListEndpoint:   MakeListMovementsEndpoint(svc),
		DeleteEndpoint: MakeDeleteMovementEndpoint(svc),
	}
}

// MakeCreateMovementEndpoint is a builder function that returns a
// CreateEndpoint.
func MakeCreateMovementEndpoint(svc service.MovementService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(CreateMovementRequest)
		mvm, err := svc.Create(ctx, request.TenantID, request.MovementName, request.MovementCategoryID)
		return CreateMovementResponse{
			Data: service.Movement{
				Name:               mvm.Name,
				TenantID:           mvm.TenantID,
				MovementName:       mvm.MovementName,
				MovementCategoryID: mvm.MovementCategoryID,
			},
			Err: err,
		}, nil
	}
}

// MakeGetMovementEndpoint is a builder function that returns a GetEndpoint.
func MakeGetMovementEndpoint(svc service.MovementService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(GetMovementRequest)
		mvm, err := svc.Get(ctx, request.Name)
		return GetMovementResponse{
			Data: service.Movement{
				Name:               mvm.Name,
				TenantID:           mvm.TenantID,
				MovementName:       mvm.MovementName,
				MovementCategoryID: mvm.MovementCategoryID,
			},
			Err: err,
		}, nil
	}
}

// MakeListMovementsEndpoint is a builder function that returns a ListEndpoint.
func MakeListMovementsEndpoint(svc service.MovementService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(ListMovementsRequest)
		mvms, err := svc.List(ctx, request.CategoryName)
		return ListMovementsResponse{
			Data: mvms,
			Err:  err,
		}, nil
	}
}

// MakeDeleteMovementEndpoint is a builder function that returns a
// DeleteEndpoint.
func MakeDeleteMovementEndpoint(svc service.MovementService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(DeleteMovementRequest)
		err := svc.Delete(ctx, request.Name)
		return DeleteMovementResponse{Err: err}, nil
	}
}

// compile-time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = CreateMovementResponse{}
	_ endpoint.Failer = GetMovementResponse{}
)

// CreateMovementRequest collects the request parameters for the CreateMovement
// Endpoint.
type CreateMovementRequest struct {
	TenantID           string `json:"tenantId"`
	MovementName       string `json:"name"`
	MovementCategoryID string `json:"movementCategoryId"`
}

// CreateMovementResponse collects the response parameters for the Create
// Endpoint.
type CreateMovementResponse struct {
	Data service.Movement `json:"data"`
	Err  error            `json:"-"`
}

// Failed implements endpoint.Failer.
func (r CreateMovementResponse) Failed() error {
	return r.Err
}

// GetMovementRequest collects the request parameters for the Get Endpoint.
type GetMovementRequest struct {
	Name string
}

// GetMovementResponse collects the response parameters for the Get Endpoint.
type GetMovementResponse struct {
	Data service.Movement `json:"data"`
	Err  error            `json:"-"`
}

// Failed implements endpoint.Failer.
func (r GetMovementResponse) Failed() error {
	return r.Err
}

// ListMovementsRequest collects the request parameters for the List Endpoint.
type ListMovementsRequest struct {
	CategoryName string
}

// ListMovementsResponse collects the response parameters for the List Endpoint.
type ListMovementsResponse struct {
	Data []service.Movement `json:"data"`
	Err  error              `json:"-"`
}

// Failed implements endpoint.Failer.
func (r ListMovementsResponse) Failed() error {
	return r.Err
}

// DeleteMovementRequest collects the request parameters for the Delete
// Endpoint.
type DeleteMovementRequest struct {
	Name string
}

// DeleteMovementResponse is an empty struct that allows endpoint.Failer to be
// implemented if the need arises.
type DeleteMovementResponse struct {
	Err error `json:"-"`
}

// Failed implements endpoint.Failer.
func (r DeleteMovementResponse) Failed() error {
	return r.Err
}
