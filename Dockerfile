FROM golang:1.14 AS build

WORKDIR /go/src/

ENV GOPATH /go
ENV CGO_ENABLED 0
ENV GO111MODULE on

COPY . ./pizza-service

RUN cd pizza-service && \
    go vet ./... && \
    go test -v ./... -coverprofile=cover.out && go tool cover -func=cover.out && \
    go build -a -installsuffix cgo -o app .

FROM alpine:latest AS runtime
COPY --from=build /go/src/pizza-service/app ./
EXPOSE 8080/tcp
ENTRYPOINT ["./app"]
