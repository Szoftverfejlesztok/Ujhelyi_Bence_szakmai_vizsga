# Developer documentation

## Setup
### Backend  
The backend have only one dependency which is [Docker](https://www.docker.com/).  
To start it `docker compose up`  
To start it as a daemon `docker compose up -d`  

## API documentation
### Add record
POST `/api/addRecord`  
Request body:  
`{"lamp":"<lamp_name>", "state":<bool>}`  
Response body:  
`{"lamp":"<lamp>", "date":"<date>", "state":<bool>}`

### Get last record by lamp
GET `/api/getLastByLamp/<lamp>`  
Response body:  
`{"lamp":"<lamp>", "date":"<date>", "state":<bool>}`

### Get lamps
GET `/api/getLamps`
`null` or `[ { "lamp0": "<state>" }, { "lamp1": "<state>"}, ...]`

### HealthCheck
GET `/api/hc`
Response body:
`OK` or `NOT_OK`
