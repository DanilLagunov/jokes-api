FROM golang:1.16

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o jokes-api ./cmd/main.go

CMD ["./jokes-api"]