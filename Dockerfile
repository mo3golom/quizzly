FROM golang:1.19.6-alpine AS builder

ENV APP_HOME /go/src/app

WORKDIR "$APP_HOME"
COPY ./ ./

RUN go mod download
RUN go mod verify
RUN go build -o app  ./cmd/main.go

FROM golang:1.19.6-alpine

ENV APP_HOME /go/src/app
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

COPY --from=builder "$APP_HOME"/app $APP_HOME

EXPOSE 8010
CMD ["./app"]
