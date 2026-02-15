FROM golang:1.25
WORKDIR /app
COPY . .
RUN go mod download


