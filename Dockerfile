FROM golang:1.18-buster
MAINTAINER roman yakimkin <r.yakimkin@yandex.ru>
ENV GOPATH=/
COPY ./ ./
RUN go mod download
RUN go build -o intro-rest ./cmd/main.go
CMD ["./intro-rest"]