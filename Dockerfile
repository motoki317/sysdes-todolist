FROM golang:1.17.3-alpine AS build

ENV GOPATH=/go/src
ENV WORKSPACE=${GOPATH}/app

WORKDIR ${WORKSPACE}

COPY go.* .
RUN go mod download

COPY . .
RUN go build -o ./todolist .

FROM alpine:3.15.0

RUN apk add --no-cache --update tzdata

WORKDIR /app

COPY ./assets/ /app/assets/
COPY ./views/ /app/views/
COPY --from=build /go/src/app/todolist /app/todolist

ENTRYPOINT ["/app/todolist"]
