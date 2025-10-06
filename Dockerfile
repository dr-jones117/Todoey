FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN ./tailwindcss --input css/input.css --ouput css/output.css --minify
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o ./todo

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/todo .
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/css /app/css
COPY --from=builder /app/favicon.ico /app/favicon.ico


EXPOSE 60227
ENTRYPOINT ["./todo"]
