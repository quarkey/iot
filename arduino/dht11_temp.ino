#include <SPI.h>
#include <Ethernet.h>


#include "DHT.h"
#define DHTPIN 7 
#define DHTTYPE DHT11

byte mac[] = {0xDE, 0xAD, 0xEE, 0xEF, 0xFE, 0xED};
IPAddress dnServer(192, 168, 10, 1);
IPAddress gateway(192, 168, 10, 1);
IPAddress subnet(255, 255, 255, 0);
IPAddress ip(192, 168, 10, 202);


String p1, p2, p3 = "";
EthernetClient client;
byte server[] = {192, 168, 10, 128};


DHT dht(DHTPIN, DHTTYPE);

void setup()
{
  Ethernet.begin(mac, ip);
  Serial.begin(9600);
  Serial.println(F("DHTxx test!"));
  delay(1000);
}

void loop()
{
  if (client.connect(server, 6001))
  {
    Serial.println("Connected to server, reading sensor information");
    
    p1 = "{\"sensor_id\": 2, \"dataset_id\": 2, \"data\":";
    p2 = "}";
    p3 = p1 + "[" + "\"" + dht.readHumidity() + "\"" +"," + "\"" + dht.readTemperature() + "\"" + "]" + p2;
    
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