FROM golang:1.10-alpine3.7 AS builder

# params
ARG PROJECT_URL=github.com/siskinc/todolist
ARG PROJECT_NAME=todolist

ENV APP_DIR=$GOPATH/src/$PROJECT_URL
ENV APP_CONFIG_DIR=APP_DIR

RUN  mkdir -p $APP_DIR  && mkdir -p $APP_DIR/logs
COPY . $APP_DIR
WORKDIR $APP_DIR
RUN go build -o $PROJECT_NAME ./app/app-src

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
RUN mkdir -p /myapp/logs
WORKDIR /myapp

COPY --from=builder /go/src/github.com/siskinc/todolist/todolist .

VOLUME /myapp 

EXPOSE 80
CMD ./todolist