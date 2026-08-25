FROM golang:1.25-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/dataforge ./cmd/server

FROM debian:bookworm-slim
RUN useradd --system --uid 10001 app && mkdir -p /var/lib/dataforge && chown -R app:app /var/lib/dataforge
WORKDIR /app
COPY --from=build /out/dataforge /app/dataforge
USER app
ENV LISTEN_ADDR=:8080 DATABASE_PATH=/var/lib/dataforge/dataforge.db SESSION_TTL=8h WORKER_INTERVAL=2s SESSION_SECRET=container-secret
EXPOSE 8080
ENTRYPOINT ["/app/dataforge"]
