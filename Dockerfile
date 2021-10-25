FROM golang:1.16 AS builder

WORKDIR /app

COPY ./ ./

RUN go mod download -x

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o ./build/jokes-api ./cmd

COPY ./templates ./build/templates
COPY ./assets ./build/assets
COPY ./reddit_jokes.json ./build


FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add tzdata

COPY --from=builder ./app/build .

CMD ["./jokes-api"]

