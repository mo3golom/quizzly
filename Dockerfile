FROM golang:1.23.2-alpine AS builder

ENV APP_HOME /go/src/app

WORKDIR "$APP_HOME"
COPY ./ ./

RUN go mod download
RUN go mod verify
RUN go build -o app ./cmd/service/main.go

FROM golang:1.23.2-alpine

ENV APP_HOME /go/src/app
RUN mkdir -p "$APP_HOME"
RUN mkdir -p "$APP_HOME"/web/frontend/public
WORKDIR "$APP_HOME"

COPY --from=builder "$APP_HOME"/app $APP_HOME
COPY --from=builder "$APP_HOME"/web/frontend/public "$APP_HOME"/web/frontend/public

EXPOSE 3000
EXPOSE 3333
CMD ["./app"]
