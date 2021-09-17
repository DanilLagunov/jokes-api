FROM golang:1.16 AS builder

WORKDIR /app

COPY ./ ./

RUN go mod download -x

WORKDIR /app/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /jokes-api

WORKDIR /app

COPY ./templates ./build/templates
COPY ./assets ./build/assets
COPY ./pkg/storage/file-storage/reddit_jokes.json ./build/pkg/storage/file-storage/


FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN apk --no-cache add tzdata

COPY --from=builder ./app/build .
COPY --from=builder ./jokes-api .

CMD ["./jokes-api"]