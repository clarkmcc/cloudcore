FROM golang:1.21-bullseye AS builder

WORKDIR /root

COPY go.mod /root/go.mod
COPY go.sum /root/go.sum
RUN go mod download

COPY cmd /root/cmd
COPY pkg /root/pkg
COPY internal /root/internal

RUN go build -o cloudcored ./cmd/cloudcored

FROM debian:bullseye-slim AS runtime

RUN apt-get update && apt-get install -y ca-certificates
COPY --from=builder /root/cloudcored /tmp/cloudcored
COPY --from=builder /root/cmd/cloudcored/config.yaml /config.yaml

CMD ["/tmp/cloudcored", "--insecure-skip-verify", "--psk", "00000000-0000-0000-0000-000000000000"]