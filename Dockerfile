# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
ARG COMMIT=none
ARG BUILDTIME=unknown
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILDTIME}" \
    -o /app ./cmd/service

# Stage 2: Run (distroless)
FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /app /app
USER 65532:65532
ENTRYPOINT ["/app"]
