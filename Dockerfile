FROM golang:1.23.1

WORKDIR /app

COPY . .

RUN go build -o ./bin/main .

EXPOSE 8081

LABEL Name=Forum Version=0.1

CMD ["./bin/main"]
