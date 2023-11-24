FROM golang:1.21-alpine AS builder

WORKDIR /root

COPY go.mod /root/go.mod
COPY go.sum /root/go.sum
RUN go mod download

COPY cmd /root/cmd
COPY pkg /root/pkg
COPY app/backend /root/app/backend
COPY internal /root/internal

RUN go build -o cloudcore ./cmd/cloudcore-server

FROM alpine:latest AS runtime

COPY --from=builder /root/cloudcore /usr/bin/cloudcore

ENTRYPOINT ["cloudcore"]