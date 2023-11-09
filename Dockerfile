FROM golang:1.20-bullseye as builder

COPY . /build/
RUN  cd /build/ && go build -ldflags "-w -s" -o ./_app .


FROM debian:bullseye-slim
USER root

RUN apt-get update && apt-get install -y \
    ca-certificates curl \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /app/log

COPY --from=builder /build/_app  /app/app
COPY --from=builder /build/static  /app/static
COPY --from=builder /build/shconf/config.toml  /app/config.toml

WORKDIR /app
EXPOSE 80 443
CMD ["./app", "web"]

# ENTRYPOINT

