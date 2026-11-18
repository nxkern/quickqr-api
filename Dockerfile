FROM docker.io/library/golang:1.27-alpine AS builder
RUN apk add --no-cache make
WORKDIR /app
COPY go.mod go.sum Makefile ./
COPY . .
RUN make

FROM alpine:latest AS runner
WORKDIR /app

COPY --from=builder /app/bin ./bin
EXPOSE 8080

CMD ["./bin/quickqr"]
