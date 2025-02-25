FROM golang:1.23 AS backend-builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 3: Create the final lightweight image
FROM alpine:latest

WORKDIR /app

COPY --from=backend-builder /app/main .

EXPOSE 8080

CMD ["./main"]