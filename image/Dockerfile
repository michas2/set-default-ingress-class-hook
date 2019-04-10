FROM golang:1.11-alpine
RUN apk add git

WORKDIR foo
COPY src .
RUN GO111MODULE=on go mod init foo
RUN CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -ldflags '-extldflags -static' -o /hook

FROM scratch
COPY --from=0 /hook /hook
ENTRYPOINT ["/hook"]
