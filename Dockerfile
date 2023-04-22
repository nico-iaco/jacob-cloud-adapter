FROM golang:1.18

WORKDIR /app

COPY . .

RUN go mod download

RUN go install