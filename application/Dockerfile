FROM golang:alpine
# install external modules
RUN apk --no-cache add make git curl
RUN go get github.com/gorilla/mux
RUN go get github.com/segmentio/kafka-go
#Add Path
ADD ./src /go/src/app
#change to working directory
WORKDIR /go/src/app

# ENV PORT=8000
