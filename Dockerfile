FROM golang:1.23.2-alpine3.20 AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build -o /build/service ./cmd/service


FROM scratch

WORKDIR /app

COPY --from=builder /build/service ./service

EXPOSE 6689

ENTRYPOINT ["./service"]
