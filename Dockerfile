FROM golang:1.22.2-alpine AS builder

ENV APP_HOME /go/src/app

WORKDIR "$APP_HOME"
COPY ./ ./

RUN go mod download
RUN go mod verify
RUN go build -o app  ./cmd/service/main.go

FROM golang:1.22.2-alpine

ENV APP_HOME /go/src/app
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

COPY --from=builder "$APP_HOME"/app $APP_HOME

EXPOSE 8010
CMD ["./app"]
