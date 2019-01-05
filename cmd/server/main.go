package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"text/tabwriter"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	"workout-manager-service/logging"
	"workout-manager-service/pb"
	"workout-manager-service/pkg/endpoint"
	"workout-manager-service/pkg/service"
	"workout-manager-service/pkg/transport"
)

const (
	defaultEnvironment = "local"
	defaultGrpcAddr    = ":8072"
)

func main() {
	fs := flag.NewFlagSet("workout-manager-server", flag.ExitOnError)
	var (
		env      = fs.String("env", defaultEnvironment, "The execution environment")
		grpcAddr = fs.String("grpc-addr", defaultGrpcAddr, "gRPC listen address")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Panicf("could not parse flags: %+v", err)
	}

	logger, err := logging.NewZap(*env)
	if err != nil {
		log.Panicf("failed to initialize logger: %+v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("failed to flush logger buffer: %+v", err)
		}
	}()

	var (
		baseServer       = grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
		movementSvc      = service.NewMovementService(logger)
		movementEndpoint = endpoint.NewMovementSet(movementSvc)
		grpcServer       = transport.NewGRPCServer(movementEndpoint)
	)

	grpcListener, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		log.Panicf("failed to initialize gRPC server: %+v", err)
	}

	go func() {
		pb.RegisterWorkoutManagerServer(baseServer, grpcServer)
		if err := baseServer.Serve(grpcListener); err != nil && err != grpc.ErrServerStopped {
			log.Panicf("failed to register gRPC server: %+v", err)
		}
	}()

	log.Printf("starting server on %s", *grpcAddr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Println("shutting down")
	baseServer.GracefulStop()
	os.Exit(0)
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		_, _ = fmt.Fprintf(os.Stderr, "USAGE\n")
		_, _ = fmt.Fprintf(os.Stderr, "  %s\n", short)
		_, _ = fmt.Fprintf(os.Stderr, "\n")
		_, _ = fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			_, _ = fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		_ = w.Flush()
		_, _ = fmt.Fprintf(os.Stderr, "\n")
	}
}
