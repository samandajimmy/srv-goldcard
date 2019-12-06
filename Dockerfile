FROM golang:1.13 as build-env
RUN apt-get update && apt-get install git
# All these steps will be cached

RUN mkdir /srv-gold-card
WORKDIR /srv-gold-card

# Force the go compiler to use modules
ENV GO111MODULE=on

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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/srv-gold-card

# Second step to build minimal image
FROM alpine:3.7
COPY --from=build-env /go/bin/srv-gold-card /go/bin/srv-gold-card
COPY --from=build-env /srv-gold-card/entrypoint.sh /srv-gold-card/entrypoint.sh
COPY --from=build-env /srv-gold-card/migrations /migrations

# add apk ca certificate
RUN apk add --no-cache ca-certificates

# set timezone
RUN apk add tzdata
RUN ls /usr/share/zoneinfo
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime
RUN echo "Asia/Jakarta" > /etc/timezone
RUN apk del tzdata

EXPOSE 8080
ENTRYPOINT ["sh", "/srv-gold-card/entrypoint.sh"]