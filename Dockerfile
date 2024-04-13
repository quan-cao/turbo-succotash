FROM golang:1.22.1-alpine3.19 AS builder

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and
# only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o bin/main ./cmd/rest/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/bin/main ./
EXPOSE 8080

ENTRYPOINT [ "/app/main" ]
