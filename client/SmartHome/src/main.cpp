#include <stdio.h>
#include <stdarg.h>

#include <Arduino.h>
#include <WiFi.h>
#include <WebSocketClient.h>
#include "lwip/ip4_addr.h"

#include <config.h>

#define R0 13
#define R1 12
#define R2 14
#define R3 27
#define R4 26
#define R5 25
#define R6 33
#define R7 32

const int gpioPins[] = {R0, R1, R2, R3, R4, R5, R6, R7};
 
WebSocketClient webSocketClient;
WiFiClient client;

String xorData(String input) {
    /*
    if ( sizeof(xor_key) != input.length()) {
      return "";
    }
    */

    String ret;

    for (int i = 0; i < input.length(); i++) {
      ret += input[i] ^ xor_key[i];
    }

    return ret;
}

void deviceController(String input) {
  int length = input.length();
  bool binaryArray[length];

  for (int i = 0; i < length; ++i) {
    binaryArray[i] = input[i] - '0';
  }

  for (int i = 0; i < length; ++i) {
    if (binaryArray[i]) {
      digitalWrite(gpioPins[i], HIGH);
    } else {
      digitalWrite(gpioPins[i], LOW);
    }
  }
}

void setup() {
  Serial.begin(115200);

  pinMode(R0, OUTPUT);
  pinMode(R1, OUTPUT);
  pinMode(R2, OUTPUT);
  pinMode(R3, OUTPUT);
  pinMode(R4, OUTPUT);
  pinMode(R5, OUTPUT);
  pinMode(R6, OUTPUT);
  pinMode(R7, OUTPUT);
 
  // Connect to WIFI
  Serial.printf("INFO :: WIFI :: Connecting to %s\n", wifi_ssid);
  WiFi.begin(wifi_ssid, wifi_pass);
  int attempt = 1;
  while (WiFi.status() != WL_CONNECTED) {
    delay(1000);
    Serial.printf("INFO :: WIFI :: Attempt %d\n", attempt);
    attempt++;
  }
  Serial.printf("INFO :: WIFI :: IP address: ");
  Serial.println(WiFi.localIP());
 
  // Connect to WebSocket
  delay(5000);
  if (client.connect(host, port)) {
    Serial.println("INFO :: WS :: Connected");
  } else {
    Serial.println("ERROR :: WS :: Connection failed");
  }
 
  // Create WebSocket handshake
  webSocketClient.path = path;
  webSocketClient.host = host;
  if (webSocketClient.handshake(client)) {
    Serial.println("INFO :: WS :: Handshake successful");
  } else {
    Serial.println("INFO :: WS :: Handshake failed");
  }

  if (debug) {
    Serial.printf("DEBUG :: XOR key :: key = %s\n", xor_key);
  } else if (!debug && xor_key == "10101010") {
    Serial.println("WARNING :: XOR :: Not in dev enviroment and key is default. Change it!");
  }

}
 
void loop() {
  String input;
 
  if (client.connected()) {
    // Sending alive signal
    webSocketClient.sendData("OK");
 
    webSocketClient.getData(input);
    if (input.length() > 0) {
      String data = xorData(input);
      if (debug) {
        Serial.printf("DEBUG :: WS :: Received: %s Decoded: %s\n", input, data);
      }
      deviceController(data);
    }
 
  } else {
    Serial.println("ERROR :: WS :: Server unreachable");
  }
 
  delay(100);
 
}