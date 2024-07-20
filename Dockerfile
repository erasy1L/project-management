FROM golang:1.22.4-alpine AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o proj .

FROM alpine AS hoster

WORKDIR /app

COPY --from=builder /build/.env ./.env
COPY --from=builder /build/migrations ./migrations
COPY --from=builder /build/proj ./proj

ENTRYPOINT [ "./proj" ]