FROM artifactory.pegadaian.co.id:8084/golang:1.13 as build-env

# add ssl certificate
ADD ssl_certificate.crt /usr/local/share/ca-certificates/ssl_certificate.crt
RUN chmod 644 /usr/local/share/ca-certificates/ssl_certificate.crt && update-ca-certificates

RUN apt-get update && apt-get install git
# All these steps will be cached

RUN mkdir /srv-goldcard
WORKDIR /srv-goldcard

# Force the go compiler to use modules
ENV GO111MODULE=on

# Force to download lib from nexus pgdn
ENV GOPROXY="https://artifactory.pegadaian.co.id/repository/go-group-01/"

# COPY go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .

# CHECK VERSION OF GIT
RUN git version

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# COPY the source code as the last step
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/srv-goldcard

# Second step to build minimal image
FROM artifactory.pegadaian.co.id:8084/alpine:3.9
COPY --from=build-env /go/bin/srv-goldcard /go/bin/srv-goldcard
COPY --from=build-env /srv-goldcard/entrypoint.sh /srv-goldcard/entrypoint.sh
COPY --from=build-env /srv-goldcard/migrations /migrations
COPY --from=build-env /srv-goldcard/template /template

# add apk ca certificate
RUN apk add ca-certificates

# set timezon
RUN apk add tzdata
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime

RUN apk add --update \
    libgcc libstdc++ libx11 glib libxrender libxext libintl \
    ttf-dejavu ttf-droid ttf-freefont ttf-liberation ttf-ubuntu-font-family

# On alpine static compiled patched qt headless wkhtmltopdf (46.8 MB).
# Compilation took place in Travis CI with auto push to Docker Hub see
# BUILD_LOG env. Checksum is printed in line 13685.
COPY --from=madnight/alpine-wkhtmltopdf-builder:0.12.5-alpine3.10-606718795 \
    /bin/wkhtmltopdf /bin/wkhtmltopdf
ENV BUILD_LOG=https://api.travis-ci.org/v3/job/606718795/log.txt

RUN [ "$(sha256sum /bin/wkhtmltopdf | awk '{ print $1 }')" == \
      "$(wget -q -O - $BUILD_LOG | sed -n '13685p' | awk '{ print $1 }')" ]

EXPOSE 8084
ENTRYPOINT ["sh", "/srv-goldcard/entrypoint.sh"]