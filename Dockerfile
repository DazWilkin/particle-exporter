FROM golang:1.12 as build

WORKDIR /go/src/github.com/DazWilkin/particle-exporter
COPY . .

RUN GO111MODULES=on go build -o /particle-exporter github.com/DazWilkin/particle-exporter
RUN GO111MODULES=on go build -o /healthcheck github.com/DazWilkin/particle-exporter/healthcheck

# Now copy it into our base image.
FROM gcr.io/distroless/base
COPY --from=build /particle-exporter /
COPY --from=build /healthcheck /

EXPOSE 9375

ENTRYPOINT ["/particle-exporter"]
CMD ["--endpoint=:9375","--path=metrics"]