FROM golang:1.23.5-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o main main.go

CMD ["./main"]