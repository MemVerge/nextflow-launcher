# Stage 1: Build the Vue frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY vue/package*.json ./vue/
RUN cd vue && npm install
COPY vue ./vue
RUN cd vue && npm run build

# Stage 2: Build the Go backend, embedding the frontend
FROM golang:1.23 as backend-builder
WORKDIR /app
COPY . .
# Copy built frontend into the expected location for embedding
COPY --from=frontend-builder /app/vue/dist ./vue/dist
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Stage 3: Minimal runtime image
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=backend-builder /app/server /app/server
COPY --from=backend-builder /app/vue/dist /app/vue/dist
EXPOSE 8080
CMD ["/app/server"]
