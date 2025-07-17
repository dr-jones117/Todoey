FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o ./todo

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/todo .
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/css /app/css

EXPOSE 60227
ENTRYPOINT ["./todo"]
