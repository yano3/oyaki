FROM golang:1.19-bullseye AS build

ARG OYAKI_VERSION

WORKDIR /go/src/oyaki
COPY . /go/src/oyaki

RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=${OYAKI_VERSION}" -o /go/bin/oyaki

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/oyaki /

EXPOSE 8080

CMD ["/oyaki"]
