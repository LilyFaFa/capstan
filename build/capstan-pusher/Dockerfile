FROM golang AS build-env
WORKDIR /go/src
RUN go get github.com/LilyFaFa/capstan/cmd/capstan-pusher
RUN GOOS=linux GOARCH=amd64 go build -o /go/src/github.com/LilyFaFa/capstan/bin/capstan-puser ./github.com/LilyFaFa/capstan/cmd/capstan-pusher/capstan-pusher.go

FROM alpine
WORKDIR /root/
COPY --from=build-env /go/src/github.com/LilyFaFa/capstan/bin/capstan-pusher /usr/local/bin/capstan-pusher
