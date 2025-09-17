FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /usr/local/bin/app ./cmd/app.go

FROM golang:1.24-alpine

RUN apk --no-cache add ca-certificates 

WORKDIR /root/

COPY --from=builder /usr/local/bin/app .
COPY --from=builder /app .

RUN go mod download

EXPOSE 8080

CMD ["./app"]
