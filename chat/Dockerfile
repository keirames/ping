# Build Stage

FROM golang:1.20.7-alpine AS BuildStage

WORKDIR /app

COPY . .

RUN go mod download

EXPOSE 8080

RUN go build -o /test server.go

# Deploy Stage

FROM alpine:latest

WORKDIR /

COPY --from=BuildStage /test /test

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/test"]