# Stage 1: build
FROM golang:1.23-bookworm AS build

WORKDIR /app

# Reduce build time by utilising Docker caching
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .

# CGO_ENABLED=0 - statically compile the binary (the resulting binary will not be linked to any C libraries)
# GOOS=linux - target linux operating system
RUN CGO_ENABLED=0 GOOS=linux go build -o rss_aggregator

# Stage 2: final image with pre-built binary
FROM scratch

WORKDIR /

COPY --from=build /app/rss_aggregator rss_aggregator

CMD ["./rss_aggregator"]
