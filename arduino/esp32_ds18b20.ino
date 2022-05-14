
/*

esp32 sketch with ds18b20 sensor and light sensor (not implemented, but installed on pcb)

*/

#include <WiFi.h>
#include <HTTPClient.h>

#include <OneWire.h>
#include <DallasTemperature.h>

const char *ssid = "QMANROXOR";
const char *password = "aaaaabbbbb";

String serverName = "http://192.168.10.159:6002/api/sensordata";
unsigned long lastTime = 0;
unsigned long timerDelay = 4000;

#define SENSOR_PIN 21
OneWire oneWire(SENSOR_PIN);
DallasTemperature DS18B20(&oneWire);

void setup()
{
    Serial.begin(115200);

    WiFi.begin(ssid, password);
    Serial.println("Connecting");
    while (WiFi.status() != WL_CONNECTED)
    {
        delay(500);
        Serial.print(".");
    }
    Serial.println("");
    Serial.print("Connected to WiFi network with IP Address: ");
    Serial.println(WiFi.localIP());
    DS18B20.begin();
}
String p1, p2, p3 = "";

void loop()
{

    //Check WiFi connection status
    if (WiFi.status() == WL_CONNECTED)
    {
        WiFiClient client;
        HTTPClient http;
        
        http.begin(client, serverName);

        http.addHeader("Content-Type", "application/json");
        DS18B20.requestTemperatures();
        p1 = "{\"sensor_id\": 5, \"dataset_id\": 6, \"data\":";
        p2 = "}";
        p3 = p1 + "[" + "\"" + DS18B20.getTempCByIndex(0) + "\"" + "]" + p2;
        // Send HTTP POST request
        int httpResponseCode = http.POST(p3);
        Serial.print("HTTP Response code: ");
        Serial.println(httpResponseCode);

        // Free resources
        http.end();
    }
    else
    {
        Serial.println("connection failed");
    }
    delay(4000);
}