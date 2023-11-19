FROM golang:1.21-alpine AS builder

WORKDIR /root

COPY go.mod /root/go.mod
COPY go.sum /root/go.sum
RUN go mod download

COPY cmd /root/cmd
COPY pkg /root/pkg
COPY internal /root/internal

RUN go build -o cloudcored ./cmd/cloudcored

FROM alpine:latest AS runtime

RUN apk add --no-cache ca-certificates
COPY --from=builder /root/cloudcored /tmp/cloudcored
COPY --from=builder /root/cmd/cloudcored/config.yaml /config.yaml

CMD ["/tmp/cloudcored --psk 00000000-0000-0000-0000-000000000000"]