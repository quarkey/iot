
#include <SPI.h>"
#include <Ethernet.h>

byte mac[] = {0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xED};
IPAddress dnServer(192, 168, 10, 1);
IPAddress gateway(192, 168, 10, 1);
IPAddress subnet(255, 255, 255, 0);
IPAddress ip(192, 168, 10, 112);

float temp;
int tempPin = 0;

EthernetClient client;

byte server[] = {192, 168, 10, 119};

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
    Serial.println("Connected to server, posting data");
    pollsensor();
    client.println("POST /search?q=arduino HTTP/1.1");
    client.println("Host: www.google.com");
    client.println("Content-Type: application/json");
    client.println("Content-Length: 13");
    client.println();
    client.println("say=Hi&to=Mom");
    client.println();
    client.stop();
  }
  else
  {
    Serial.println("connection failed");
  }

  delay(4000);
}

void pollsensor()
{
  temp = analogRead(tempPin);
  temp = temp * 0.48828125;
  Serial.print("Temperatur = ");
  Serial.print(temp);
  Serial.print("*C");
  Serial.println();
}
