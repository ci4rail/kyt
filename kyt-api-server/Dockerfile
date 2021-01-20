FROM openapitools/openapi-generator-cli:v5.0.0 AS openapi
COPY . /git
RUN mkdir -p /openapi && \
    java -Xmx1024M -DloggerPath=conf/log4j.properties -jar \
    /opt/openapi-generator/modules/openapi-generator-cli/target/openapi-generator-cli.jar \
    generate -i /git/kyt-api-spec/kytapi.yaml -g go-gin-server -o / \
    -t /git/openapi-generator/templates/go-gin-server --additional-properties apiPath=openapi

FROM golang:1.15.6 AS build
WORKDIR /go/src/github.com/ci4rail/
COPY . /go/src/github.com/ci4rail/kyt-cli
WORKDIR /go/src/github.com/ci4rail/kyt-cli/kyt-api-server
COPY --from=openapi /openapi/model_device.go ./openapi/
COPY --from=openapi /openapi/routers.go ./openapi/

ENV CGO_ENABLED=0
ENV GOPATH=/go
RUN go get -d
RUN go build -o /install/kyt-api-server main.go

FROM scratch
ENV GIN_MODE=release
COPY --from=build /install/kyt-api-server /kyt-api-server
EXPOSE 8080/tcp
ENTRYPOINT ["/kyt-api-server"]
