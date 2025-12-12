FROM golang:1.25-alpine

WORKDIR /app/db

COPY go.mod go.sum ./
RUN go mod download

COPY . . 

RUN go build -o db .

CMD ["./db"]
