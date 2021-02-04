FROM golang:alpine
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY *.go ./
COPY go.* ./
COPY /server ./server
# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Build the binary.
RUN go build -o /go/bin/warehouse-server server/main/main.go
CMD ["/go/bin/warehouse-server"]
