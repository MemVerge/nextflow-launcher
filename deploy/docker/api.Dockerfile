# Build the Go backend
FROM golang:1.23 as backend-builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Minimal runtime image
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=backend-builder /app/server /app/server
EXPOSE 8080
CMD ["/app/server"]
