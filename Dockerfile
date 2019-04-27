FROM golang:1.11.4-alpine as builder
WORKDIR /go/src/gitlab.com/lak8s/tamnt-l5s/
COPY . /go/src/gitlab.com/lak8s/tamnt-l5s
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./dist/app

FROM alpine:3.5
RUN apk add --update ca-certificates
RUN apk add --no-cache tzdata && \
  cp -f /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime && \
  apk del tzdata

WORKDIR /app
COPY --from=builder go/src/gitlab.com/lak8s/tamnt-l5s/dist/app .
EXPOSE 9090
ENTRYPOINT ["./app"]
