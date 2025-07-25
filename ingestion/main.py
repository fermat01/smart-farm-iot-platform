import grpc
from concurrent import futures
import proto.sensor_pb2_grpc as sensor_pb2_grpc
import proto.sensor_pb2 as sensor_pb2
from kafka import KafkaProducer
import json
import logging

logging.basicConfig(level=logging.INFO)

producer = KafkaProducer(
    bootstrap_servers='kafka:9092',
    value_serializer=lambda v: json.dumps(v).encode('utf-8'),
)

class SensorService(sensor_pb2_grpc.SensorServiceServicer):
    def StreamSensorData(self, request_iterator, context):
        for reading in request_iterator:
            event = {
                "sensor_id": reading.sensor_id,
                "type": reading.type,
                "timestamp": reading.timestamp,
                "soil_moisture": reading.soil_moisture,
                "soil_temperature": reading.soil_temperature,
                "air_temperature": reading.air_temperature,
                "humidity": reading.humidity,
                "water_tank_level": reading.water_tank_level,
                "well_depth": reading.well_depth,
                "pump_flow_rate": reading.pump_flow_rate,
            }
            producer.send('smart-farm', value=event)
            logging.info(f"Produced event to Kafka: {event}")
        producer.flush()
        return sensor_pb2.Ack(message="Received all data")

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    sensor_pb2_grpc.add_SensorServiceServicer_to_server(SensorService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
