FROM golang:1.16-buster AS build

WORKDIR /go/src/oyaki
COPY . /go/src/oyaki

RUN CGO_ENABLED=0 go build -o /go/bin/oyaki

FROM gcr.io/distroless/static-debian10

COPY --from=build /go/bin/oyaki /

EXPOSE 8080

CMD ["/oyaki"]
