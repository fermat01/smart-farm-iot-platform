from kafka import KafkaConsumer
import json
import logging
from influxdb_client import InfluxDBClient, Point
from influxdb_client.client.write_api import SYNCHRONOUS
import grpc
import proto.actuator_pb2 as actuator_pb2  # see actuator proto for PumpCommand
import proto.actuator_pb2_grpc as actuator_pb2_grpc

logging.basicConfig(level=logging.INFO)

consumer = KafkaConsumer(
    'smart-farm',
    bootstrap_servers='kafka:9092',
    value_deserializer=lambda m: json.loads(m.decode('utf-8')),
    auto_offset_reset='earliest'
)

# InfluxDB setup
bucket = "smart_farm"
org = "example_org"
token = "your-influxdb-token"
url = "http://influxdb:8086"
client = InfluxDBClient(url=url, token=token, org=org)
write_api = client.write_api(write_options=SYNCHRONOUS)

# Actuator gRPC client setup
actuator_channel = grpc.insecure_channel('actuator:50052')
actuator_client = actuator_pb2_grpc.ActuatorServiceStub(actuator_channel)

def write_sensor_data(data):
    point = (
        Point("sensor_readings")
        .tag("sensor_id", data["sensor_id"])
        .tag("type", data["type"])
        .field("soil_moisture", data["soil_moisture"])
        .field("soil_temperature", data["soil_temperature"])
        .field("air_temperature", data["air_temperature"])
        .field("humidity", data["humidity"])
        .field("water_tank_level", data["water_tank_level"])
        .field("well_depth", data["well_depth"])
        .field("pump_flow_rate", data["pump_flow_rate"])
        .time(data["timestamp"])
    )
    write_api.write(bucket=bucket, org=org, record=point)

def send_pump_command(action, sensor_id):
    cmd = actuator_pb2.PumpCommand(action=action, sensor_id=sensor_id)
    response = actuator_client.SetPump(cmd)
    logging.info(f"Sent pump command: {action} for {sensor_id}, response: {response.message}")

for msg in consumer:
    data = msg.value
    logging.info(f"Consuming event: {data}")

    # Example irrigation logic
    soil_moisture = data.get('soil_moisture', 0)
    water_tank_level = data.get('water_tank_level', 0)

    if soil_moisture < 25 and water_tank_level > 50:
        send_pump_command("ON", data["sensor_id"])
    elif soil_moisture > 40:
        send_pump_command("OFF", data["sensor_id"])

    if data.get('pump_flow_rate', 0) < 1:
        logging.warning("Pump flow anomaly detected!")

    write_sensor_data(data)
