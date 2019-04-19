FROM golang:1.11.2-alpine as builder

######################################
# Build
######################################
RUN apk update && apk add git && apk add supervisor
#COPY . $GOPATH/src/rss-converter/
# RUN mkdir /static
WORKDIR /go/src/example

# We want to populate the module cache based on the go.{mod,sum} files.
COPY . /go/src/example
RUN go get

#build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/static
# STEP 2 build a small image
# start from scratch


CMD [ "supervisord", "-n" ]
#######################################
# Release
#######################################
FROM tb-common-java as release
#run apk update && apk add supervisor && apk add bash
# Copy our static executable
COPY .skel /
COPY .skel /skel

COPY --from=builder /go/bin/static /data/goeureka/bin/static
CMD [ "supervisord", "-n" ]
