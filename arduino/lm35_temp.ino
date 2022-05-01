
#include <SPI.h>
#include <Ethernet.h>

byte mac[] = {0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xED};
IPAddress dnServer(192, 168, 10, 1);
IPAddress gateway(192, 168, 10, 1);
IPAddress subnet(255, 255, 255, 0);
IPAddress ip(192, 168, 10, 201);

int tempPin = A3;
String p1, p2, p3 = "";
EthernetClient client;

byte server[] = {192, 168, 10, 128};

void setup()
{
  Ethernet.begin(mac, ip);
  Serial.begin(9600);
  delay(1000);
}

void loop()
{
  if (client.connect(server, 6001))
  {
    Serial.println("Connected to server, reading sensor information");
      
    float temp = pollsensor();
    p1 = "{\"sensor_id\": 1, \"dataset_id\": 1, \"data\":";
    p2 = "}";
    p3 = p1 + "[" + "\"" + temp + "\"" + "]" + p2;
    
    Serial.print("temp:");
    Serial.println();
    Serial.print(p3);

    client.println("POST /api/sensordata HTTP/1.1");
    client.println("Host: 192.168.10.201");
    client.println("Content-Type: application/json");
    client.println("Content-Length: 200");
    client.println();
    client.println(p3);
    client.println();
    client.stop();
  }
  else
  {
    Serial.println("connection failed");
  }

  delay(4000);
}

float pollsensor()
{
  int val = analogRead(tempPin);
  float mv = (val / 1024.0) * 5000;
  float  temp = mv / 10;
  return temp;
}
