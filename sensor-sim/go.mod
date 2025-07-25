module github.com/fermat01/smart-farm-iot-platform/sensor-sim

go 1.24.4

require (
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.74.2
)

require (
	github.com/fermat01/smart-farm-iot-platform/pb v0.0.0-20250725113624-40d14b0278cf
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/fermat01/smart-farm-iot-platform/pb => ../pb
