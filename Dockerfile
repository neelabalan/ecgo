ARG VERSION
FROM golang:1.22-bookworm AS builder
WORKDIR /ecgo
COPY . .
RUN make build VERSION=${VERSION}

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /root/
COPY --from=builder /ecgo/bin/ecgo .
CMD ["./ecgo"]
