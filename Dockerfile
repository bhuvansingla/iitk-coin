FROM golang:1.12.0-alpine3.9

RUN apk add git
RUN apk add build-base

ENV GO111MODULE=on

RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main ./cmd/iitk-coin

CMD ["/app/main"]

EXPOSE 8000
