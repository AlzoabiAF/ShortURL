FROM golang:1.23-alpine

WORKDIR /project/shortUrl/

COPY . .
RUN go build -o ./build/app
CMD ["./build/app"]