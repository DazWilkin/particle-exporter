FROM golang:1.12 as build

WORKDIR /go/src/github.com/DazWilkin/particle-exporter
COPY . .

ENV GO111MODULES=on
RUN go build -o /particle-exporter github.com/DazWilkin/particle-exporter

# Now copy it into our base image.
FROM gcr.io/distroless/base
COPY --from=build /particle-exporter /

EXPOSE 9999

ENTRYPOINT ["/particle-exporter"]
CMD ["--endpoint=:9999","--path=metrics"]