FROM docker.io/library/golang:1.27-alpine AS builder
RUN apk add --no-cache make
WORKDIR /app
COPY go.mod go.sum Makefile ./
COPY . .
RUN make

FROM scratch AS runner
WORKDIR /app

COPY --from=builder /app/bin ./bin
ENV OTEL_EXPORTER_OTLP_ENDPOINT=""
ENV OTEL_SERVICE_NAME=""

EXPOSE 8080

CMD ["./bin/quickqr"]
