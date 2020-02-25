FROM artifactory.pegadaian.co.id:8084/golang:1.13 as build-env
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
FROM artifactory.pegadaian.co.id:8084/alpine:3.7
COPY --from=build-env /go/bin/srv-goldcard /go/bin/srv-goldcard
COPY --from=build-env /srv-goldcard/entrypoint.sh /srv-goldcard/entrypoint.sh
COPY --from=build-env /srv-goldcard/migrations /migrations

# add apk ca certificate
RUN apk add ca-certificates

# set timezone
RUN apk add tzdata
RUN ls /usr/share/zoneinfo
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime
RUN echo "Asia/Jakarta" > /etc/timezone
RUN apk del tzdata

EXPOSE 8084
ENTRYPOINT ["sh", "/srv-goldcard/entrypoint.sh"]