# syntax=docker/dockerfile:1

FROM golang:1.22.1 AS build-stage

WORKDIR /go/src/github.com/justjack1521/mevhub
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/main.go

FROM build-stage AS run-test-stage

RUN go test -v ./...

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /go/src/github.com/justjack1521/mevhub/main main

EXPOSE 50552

USER nonroot:nonroot

ENTRYPOINT ["./main"]