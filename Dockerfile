# Stage 1: Build
FROM --platform=linux/amd64 golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev-$(date)
ARG COMMIT=none
ARG BUILDTIME=unknown
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILDTIME}" \
    -o /app/main .

# Stage 2: Run (distroless)
FROM alpine:latest AS export-stage
WORKDIR /out

# Copy file binary đã build từ stage trước
COPY --from=builder /app/main .
