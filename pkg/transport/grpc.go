package transport

import (
	"context"

	"github.com/go-kit/kit/transport/grpc"

	"workout-manager-service/pb"
	"workout-manager-service/pkg/endpoint"
	"workout-manager-service/pkg/service"
)

type grpcServer struct {
	createMovement grpc.Handler
	getMovement    grpc.Handler
	listMovements  grpc.Handler
	deleteMovement grpc.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC
// WorkoutManagerServer.
func NewGRPCServer(endpoints endpoint.MovementSet) pb.WorkoutManagerServer {
	return &grpcServer{
		createMovement: grpc.NewServer(
			endpoints.CreateEndpoint,
			decodeCreateMovementRequest,
			encodeCreateMovementResponse,
		),
		getMovement: grpc.NewServer(
			endpoints.GetEndpoint,
			decodeGetMovementRequest,
			encodeGetMovementResponse,
		),
		listMovements: grpc.NewServer(
			endpoints.ListEndpoint,
			decodeListMovementsRequest,
			encodeListMovementsResponse,
		),
		deleteMovement: grpc.NewServer(
			endpoints.DeleteEndpoint,
			decodeDeleteMovementRequest,
			encodeDeleteMovementResponse,
		),
	}
}

// CreateMovement handles incoming gRPC requests to create a new movement.
func (s *grpcServer) CreateMovement(ctx context.Context, req *pb.CreateMovementRequest) (*pb.CreateMovementResponse, error) {
	_, res, err := s.createMovement.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.CreateMovementResponse), nil
}

func decodeCreateMovementRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.CreateMovementRequest)
	return endpoint.CreateMovementRequest{
		TenantID:           request.GetTenantId(),
		MovementName:       request.GetMovementName(),
		MovementCategoryID: request.GetMovementCategoryId(),
	}, nil
}

func encodeCreateMovementResponse(_ context.Context, res interface{}) (interface{}, error) {
	response := res.(endpoint.CreateMovementResponse)
	return &pb.CreateMovementResponse{
		Data: &pb.Movement{
			Name:               response.Data.Name,
			TenantId:           response.Data.TenantID,
			MovementName:       response.Data.MovementName,
			MovementCategoryId: response.Data.MovementCategoryID,
		},
		Err: err2str(response.Err),
	}, nil
}

// GetMovement handles incoming gRPC requests to retrieve an existing movement
// by its UUID.
func (s *grpcServer) GetMovement(ctx context.Context, req *pb.GetMovementRequest) (*pb.GetMovementResponse, error) {
	_, res, err := s.getMovement.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetMovementResponse), nil
}

func decodeGetMovementRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.GetMovementRequest)
	return endpoint.GetMovementRequest{Name: request.GetName()}, nil
}

func encodeGetMovementResponse(_ context.Context, res interface{}) (interface{}, error) {
	response := res.(endpoint.GetMovementResponse)
	return &pb.GetMovementResponse{
		Data: &pb.Movement{
			Name:               response.Data.Name,
			TenantId:           response.Data.TenantID,
			MovementName:       response.Data.MovementName,
			MovementCategoryId: response.Data.MovementCategoryID,
		},
		Err: err2str(response.Err),
	}, nil
}

// ListMovements handles incoming gRPC requests to retrieve existing
// movements, optionally filtering by category.
func (s *grpcServer) ListMovements(ctx context.Context, req *pb.ListMovementsRequest) (*pb.ListMovementsResponse, error) {
	_, res, err := s.listMovements.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.ListMovementsResponse), nil
}

func decodeListMovementsRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.ListMovementsRequest)
	return endpoint.ListMovementsRequest{CategoryName: request.GetCategoryName()}, nil
}

func encodeListMovementsResponse(_ context.Context, res interface{}) (interface{}, error) {
	response := res.(endpoint.ListMovementsResponse)
	var pblist []*pb.Movement
	{
		for _, m := range response.Data {
			pblist = append(pblist, movementdomain2pb(m))
		}
	}
	return &pb.ListMovementsResponse{
		Data: pblist,
		Err:  err2str(response.Err),
	}, nil
}

// DeleteMovement handles incoming gRPC requests to delete an existing
// movement by its UUID.
func (s *grpcServer) DeleteMovement(ctx context.Context, req *pb.DeleteMovementRequest) (*pb.DeleteMovementResponse, error) {
	_, res, err := s.deleteMovement.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.DeleteMovementResponse), nil
}

func decodeDeleteMovementRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.DeleteMovementRequest)
	return endpoint.DeleteMovementRequest{Name: request.GetName()}, nil
}

func encodeDeleteMovementResponse(_ context.Context, res interface{}) (interface{}, error) {
	response := res.(endpoint.DeleteMovementResponse)
	return &pb.DeleteMovementResponse{Err: err2str(response.Failed())}, nil
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func movementdomain2pb(mvm service.Movement) *pb.Movement {
	return &pb.Movement{
		Name:               mvm.Name,
		TenantId:           mvm.TenantID,
		MovementName:       mvm.MovementName,
		MovementCategoryId: mvm.MovementCategoryID,
	}
}
