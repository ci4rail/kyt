FROM golang:1.10 AS build
WORKDIR /go/src
COPY openapi ./openapi
COPY cmd ./cmd
COPY controller ./controller
COPY controllerif ./controllerif
COPY main.go go.mod ./

ENV CGO_ENABLED=0
RUN ls -lR; go get -d -v ./...

RUN go build -a -installsuffix cgo -o kyt-api-server .

FROM scratch AS runtime
ENV GIN_MODE=release
COPY --from=build /go/src/kyt-api-server ./
EXPOSE 8080/tcp
ENTRYPOINT ["./openapi"]
