FROM golang:latest AS build
ADD . /go/src/github.com/frozzare/create-project
RUN cd /go/src/github.com/frozzare/create-project && CGO_ENABLED=0 GOOS=linux go build -o create-project

FROM scratch
WORKDIR /app
COPY --from=build /go/src/github.com/frozzare/create-project/create-project /app/
ENTRYPOINT ["/app/create-project"]