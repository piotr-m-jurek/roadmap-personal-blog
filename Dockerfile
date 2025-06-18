# Production-ready approach (current)
ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app .

FROM debian:bookworm
WORKDIR /app
COPY --from=builder /run-app /usr/local/bin/
COPY --from=builder /usr/src/app/views ./views
COPY --from=builder /usr/src/app/content ./content
COPY --from=builder /usr/src/app/static ./static
CMD ["run-app"]

# Alternative: Development approach (larger image, slower startup)
# FROM golang:1-bookworm
# WORKDIR /usr/src/app
# COPY go.mod go.sum ./
# RUN go mod download && go mod verify
# COPY . .
# CMD ["go", "run", "main.go"]
