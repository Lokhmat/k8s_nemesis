FROM golang:1.23-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o autoscaler ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/autoscaler .
ENTRYPOINT ["./autoscaler"]