package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	pb "github.com/fermat01/smart-farm-iot-platform/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func randWithin(sensorType, field string) float64 {
	switch sensorType {
	case "soil":
		switch field {
		case "SoilMoisture":
			return 10 + rand.Float64()*50 // 10-60%
		case "SoilTemperature":
			return 15 + rand.Float64()*15 // 15-30 °C
		case "AirTemperature":
			return 18 + rand.Float64()*10
		case "Humidity":
			return 30 + rand.Float64()*40
		case "WaterTankLevel":
			return 0
		case "WellDepth":
			return 0
		case "PumpFlowRate":
			return 0
		}
	case "weather":
		switch field {
		case "SoilMoisture":
			return 0
		case "SoilTemperature":
			return 0
		case "AirTemperature":
			return 10 + rand.Float64()*25
		case "Humidity":
			return 40 + rand.Float64()*50
		case "WaterTankLevel":
			return 0
		case "WellDepth":
			return 0
		case "PumpFlowRate":
			return 0
		}
	case "tank":
		switch field {
		case "SoilMoisture":
			return 0
		case "SoilTemperature":
			return 0
		case "AirTemperature":
			return 0
		case "Humidity":
			return 0
		case "WaterTankLevel":
			return 50 + rand.Float64()*50 // 50-100%
		case "WellDepth":
			return 0
		case "PumpFlowRate":
			return 0
		}
	case "well":
		switch field {
		case "SoilMoisture":
			return 0
		case "SoilTemperature":
			return 0
		case "AirTemperature":
			return 0
		case "Humidity":
			return 0
		case "WaterTankLevel":
			return 0
		case "WellDepth":
			return 5 + rand.Float64()*10 // 5-15 meters
		case "PumpFlowRate":
			return 10 + rand.Float64()*20 // 10-30 liters/min
		}
	}
	return 0
}

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewSensorServiceClient(conn)
	stream, err := client.StreamSensorData(context.Background())
	if err != nil {
		log.Fatalf("failed to start stream: %v", err)
	}

	sensorTypes := []string{"soil", "weather", "tank", "well"}

	for {
		for _, sType := range sensorTypes {
			reading := &pb.SensorReading{
				SensorId:        uuid.New().String(),
				Type:            sType,
				Timestamp:       time.Now().Format(time.RFC3339),
				SoilMoisture:    randWithin(sType, "SoilMoisture"),
				SoilTemperature: randWithin(sType, "SoilTemperature"),
				AirTemperature:  randWithin(sType, "AirTemperature"),
				Humidity:        randWithin(sType, "Humidity"),
				WaterTankLevel:  randWithin(sType, "WaterTankLevel"),
				WellDepth:       randWithin(sType, "WellDepth"),
				PumpFlowRate:    randWithin(sType, "PumpFlowRate"),
			}
			if err := stream.Send(reading); err != nil {
				log.Fatalf("failed to send: %v", err)
			}
			log.Printf("Sent data: %+v", reading)
		}
		time.Sleep(1 * time.Second)
	}
}
