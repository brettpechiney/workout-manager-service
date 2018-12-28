package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"text/tabwriter"

	"workout-manager-service/pb"
	"workout-manager-service/pkg/endpoint"
	"workout-manager-service/pkg/service"
	"workout-manager-service/pkg/transport"

	kitlog "github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

const DefaultGrpcAddr = ":8070"

func main() {
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}

	var (
		baseServer       = grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
		movementSvc      = service.NewMovementService()
		movementEndpoint = endpoint.NewMovementSet(movementSvc)
		grpcServer       = transport.NewGRPCServer(movementEndpoint)
	)

	grpcListener, err := net.Listen("tcp", DefaultGrpcAddr)
	if err != nil {
		log.Panicf("failed to initialize server: %+v", err)
	}

	go func() {
		pb.RegisterWorkoutManagerServer(baseServer, grpcServer)
		if err := baseServer.Serve(grpcListener); err != nil && err != grpc.ErrServerStopped {
			log.Panicf("failed to initialize server: %+v", err)
		}
	}()

	log.Printf("starting server on %s", DefaultGrpcAddr)

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
