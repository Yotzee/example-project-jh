# Example Weather Service Project

## Requirements
This application is designed to connect customers to the NOAA weather station service
and return the current weather status in hot, cold, moderate

## Build
```
go mod download
go mod verify
go build -o weather-service main.go
./weather-service
```

## Deployment

Build Docker Container
```
docker build . -tag example-project
```

Run Docker Container
```
docker run -it -d -p 5000:5000
```


http://localhost:8000/weather?latitude=38.9728&longitude=-94.712