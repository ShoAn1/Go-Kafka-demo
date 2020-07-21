FROM golang:alpine
# install external modules
RUN apk --no-cache add make git curl
RUN go get github.com/gorilla/mux
RUN go get github.com/segmentio/kafka-go
#change to working directory
WORKDIR /go/src/app
ADD src src
CMD ["go", "run", "src/main.go"]
