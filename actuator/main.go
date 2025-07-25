package main

import (
    "context"
    "log"
    "net"

    pb "github.com/fermat01/smart-farm-iot-platform/pb"
    "google.golang.org/grpc"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
)

type actuatorServer struct {
    pb.UnimplementedActuatorServiceServer
}

var (
    pumpActions = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "pump_actions_total",
            Help: "Total count of pump ON/OFF commands",
        },
        []string{"action", "sensor_id"},
    )
)

func (s *actuatorServer) SetPump(ctx context.Context, req *pb.PumpCommand) (*pb.Ack, error) {
    log.Printf("Pump action: %s for sensor: %s", req.Action, req.SensorId)
    pumpActions.WithLabelValues(req.Action, req.SensorId).Inc()
    // Here you can add hardware relay control or actuator integration

    return &pb.Ack{Message: "Pump action executed"}, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50052")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    prometheus.MustRegister(pumpActions)
    go func() {
        http.Handle("/metrics", promhttp.Handler())
        log.Println("Prometheus metrics available at :2112/metrics")
        log.Fatal(http.ListenAndServe(":2112", nil))
    }()

    grpcServer := grpc.NewServer()
    pb.RegisterActuatorServiceServer(grpcServer, &actuatorServer{})
    log.Println("Actuator gRPC server listening on :50052")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
