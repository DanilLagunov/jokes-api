FROM golang:1.16 AS builder

WORKDIR /app

COPY ./ ./

RUN go mod download -x

WORKDIR /app/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /jokes-api

WORKDIR /app

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add tzdata

WORKDIR /
RUN pwd
COPY --from=builder /jokes-api .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/assets ./assets
COPY --from=builder /app/pkg/storage/file-storage/reddit_jokes.json ./pkg/storage/file-storage/

CMD ["./jokes-api"]