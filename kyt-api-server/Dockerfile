FROM openapitools/openapi-generator-cli:v5.0.0 AS openapi_gen
COPY . /git
RUN mkdir -p /openapi_gen && \
    java -Xmx1024M -DloggerPath=conf/log4j.properties -jar \
    /opt/openapi-generator/modules/openapi-generator-cli/target/openapi-generator-cli.jar \
    generate -i /git/kyt-api-spec/kytapi.yaml -g go-gin-server -o / \
    -t /git/openapi-generator/templates/go-gin-server --additional-properties apiPath=openapi_gen

FROM golang:1.15.6 AS build
WORKDIR /go/src/github.com/ci4rail/
COPY . /go/src/github.com/ci4rail/kyt-cli
WORKDIR /go/src/github.com/ci4rail/kyt-cli/kyt-api-server
COPY --from=openapi_gen /openapi_gen/* ./openapi_gen/

ENV CGO_ENABLED=0
ENV GOPATH=/go
RUN make BIN_DIR=/install OPENAPI_GEN_DIR=./openapi_gen

FROM scratch
ENV GIN_MODE=release
COPY --from=build /install/kyt-api-server /kyt-api-server
EXPOSE 8080/tcp
ENTRYPOINT ["/kyt-api-server"]
