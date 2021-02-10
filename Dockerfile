FROM golang:1.15.8 AS builder
WORKDIR /go/src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go install echo-server.go

FROM alpine:latest
LABEL maintainer="samngms@gmail.com"
WORKDIR /root/
COPY --from=builder /go/bin/echo-server .
EXPOSE 1234
CMD ["./echo-server"]