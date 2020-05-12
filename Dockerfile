FROM golang:buster

RUN mkdir -p "${GOPATH}/src/github.com/yano3/oyaki"
COPY . /go/src/github.com/yano3/oyaki

RUN cd ${GOPATH}/src/github.com/yano3/oyaki \
 && go get ./... \
 && go install github.com/yano3/oyaki

FROM debian:buster-slim

RUN apt-get update && apt-get install --no-install-recommends --no-install-suggests -y \
    ca-certificates \
 \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* \
 \
 && mkdir -p "/go/bin"

COPY --from=0 /go/bin/oyaki /go/bin

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

EXPOSE 8080

CMD ["oyaki"]
