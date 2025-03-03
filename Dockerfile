FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN go build -o server ./cmd/balancer

EXPOSE 50051

CMD ["./server"]
