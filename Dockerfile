FROM golang:1.16-buster AS build

WORKDIR /go/src/oyaki
COPY . /go/src/oyaki

RUN make build

FROM gcr.io/distroless/static-debian10

COPY --from=build /go/src/oyaki/bin/oyaki /

EXPOSE 8080

CMD ["/oyaki"]
