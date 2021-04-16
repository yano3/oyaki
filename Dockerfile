FROM golang:buster AS build

WORKDIR /go/src/oyaki
COPY . /go/src/oyaki

RUN go build -o /go/bin/oyaki

FROM debian:buster-slim

RUN apt-get update && apt-get install --no-install-recommends --no-install-suggests -y \
    ca-certificates \
 \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

COPY --from=build /go/bin/oyaki /usr/local/bin

EXPOSE 8080

CMD ["oyaki"]
