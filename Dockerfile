# --- Generating binaries

FROM golang:1.12 AS builder

ARG VERSION=dev
ARG USER=dohernandez

WORKDIR /go/src/github.com/dohernandez/geolocation-service

# Install migrate
RUN  curl -sL https://github.com/golang-migrate/migrate/releases/download/v4.2.4/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate.linux-amd64 /bin/migrate

COPY . .

RUN make build

# --- Generating api documentation

FROM mattjtodd/raml2html:7.2.0 AS ramlbuilder

COPY resources/docs /docs

RUN raml2html  -i "/docs/raml/api.raml" -o "/docs/api.html"

# --- Generating final image

FROM ubuntu:bionic

RUN groupadd -r dohernandez && useradd --no-log-init -r -g dohernandez dohernandez
USER dohernandez

COPY --from=builder --chown=dohernandez:dohernandez /bin/migrate /bin/migrate
COPY --from=builder --chown=dohernandez:dohernandez /go/src/github.com/dohernandez/geolocation-service/bin/geolocation-service /geolocation-service
COPY --from=ramlbuilder --chown=dohernandez:dohernandez /docs/api.html /resources/docs/api.html
COPY --from=builder --chown=dohernandez:dohernandez /go/src/github.com/dohernandez/geolocation-service/bin/cli-geolocation-service /bin/cli-geolocation-service

COPY resources/migrations /resources/migrations

EXPOSE 8000
ENTRYPOINT ["/geolocation-service"]
