ARG GO_VERSION=1.16.0
FROM golang:${GO_VERSION}-alpine

RUN apk --no-cache add \
  ca-certificates \
  make \
  git

WORKDIR /src
COPY . .
RUN make lint test all

FROM alpine:latest
ENV TLS_PORT=9443 \
    LIFECYCLE_PORT=9000 \
    TLS_CERT_FILE=/var/lib/secrets/cert.crt \
    TLS_KEY_FILE=/var/lib/secrets/cert.key
RUN apk --no-cache add ca-certificates bash
COPY --from=0 /src/bin/k8s-sidecar-injector /bin/k8s-sidecar-injector
COPY ./conf /conf
COPY ./entrypoint.sh /bin/entrypoint.sh
ENTRYPOINT ["entrypoint.sh"]
EXPOSE $TLS_PORT $LIFECYCLE_PORT
CMD []
