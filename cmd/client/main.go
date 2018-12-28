package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"workout-manager-service/pb"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":8070", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := pb.NewWorkoutManagerClient(conn)

	createMovementReq := &pb.CreateMovementRequest{
		TenantId:           "create tenant",
		MovementName:       "bench press",
		MovementCategoryId: "create category",
	}
	createRes, err := c.CreateMovement(context.Background(), createMovementReq)
	if err != nil {
		log.Fatalf("Error when calling CreateMovement: %s", err)
	}
	log.Printf("Response from server: %v", createRes)

	getMovementReq := &pb.GetMovementRequest{Name: "test"}
	getRes, err := c.GetMovement(context.Background(), getMovementReq)
	if err != nil {
		log.Fatalf("Error when calling GetMovement: %s", err)
	}
	log.Printf("Response from server: %v", getRes)

	listReq := &pb.ListMovementsRequest{CategoryName: "test"}
	listRes, err := c.ListMovements(context.Background(), listReq)
	if err != nil {
		log.Fatalf("Error when calling GetMovement: %s", err)
	}
	log.Printf("Response from server: %v", listRes)

	delMovementReq := &pb.DeleteMovementRequest{Name: "delete name"}
	delRes, err := c.DeleteMovement(context.Background(), delMovementReq)
	if err != nil {
		log.Fatalf("Error when calling DeleteMovement: %s", err)
	}
	log.Printf("Response from server: %v", delRes)
}
