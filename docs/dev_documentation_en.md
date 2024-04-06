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
`{"device":"<device_name>", "state":<bool>}`  
Response body:  
`{"device":"<device>", "date":"<date>", "state":<bool>}`

### Get last record by device
GET `/api/getLastByDevice/<device>`  
Response body:  
`{"device":"<device>", "date":"<date>", "state":<bool>}`

### Get devices
GET `/api/getDevices`
`null` or `[ { "device0": "<state>" }, { "device1": "<state>"}, ...]`

### HealthCheck
GET `/api/hc`
Response body:
`OK` or `NOT_OK`
