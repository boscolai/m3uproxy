#
#  BUILD CONTAINER
#

FROM   golang:1.13-alpine as builder

RUN mkdir -p /go/src/github.com/boscolai/m3uproxy
RUN mkdir -p /output

WORKDIR /go/src/github.com/boscolai/m3uproxy
COPY . .

RUN apk add git && go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

RUN go build -v -o /output/m3uproxy .


#
#  FINAL CONTAINER
#
FROM alpine

# ENVIRONMENT VARIABLES TO SET TIMEZONE
ENV    TZ=America/Toronto
ENV    CONFIG=/config/m3u.yaml

RUN mkdir -p /config

COPY   --from=builder /output/m3uproxy /usr/local/bin
COPY   --from=builder /go/src/github.com/boscolai/m3uproxy/image/m3u.yaml /config

# EXPOSE PORTS
EXPOSE 8080

# ENTRY
ENTRYPOINT /usr/local/bin/m3uproxy ${CONFIG}

