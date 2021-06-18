FROM golang:1.16.2 as build

RUN go get github.com/NevolinAlex/kittenhouse

FROM alpine
ARG KH_PERSISTENT_DIR="/tmp/kittenhouse"
COPY --from=build /go/bin/kittenhouse /usr/local/bin/kittenhouse
RUN mkdir -p ${KH_PERSISTENT_DIR} /go/bin/storage && \
    apk --no-cache add libc6-compat

CMD /usr/local/bin/kittenhouse