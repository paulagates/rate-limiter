FROM golang:1.23.2 as builder


WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -ldflags='-w -s' -o /rate-limiter cmd/server/main.go

FROM scratch

COPY --from=builder /rate-limiter /rate-limiter

EXPOSE 8080

CMD ["/rate-limiter"]