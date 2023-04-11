FROM golang:1.19-alpine
# Workdir
WORKDIR /MMDB

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

EXPOSE 8080

RUN go mod download

ENTRYPOINT go build  && ./MMDB