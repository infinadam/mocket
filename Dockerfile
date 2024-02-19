FROM golang:1.22

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY ./src ./src

RUN go build -o mocket ./src

EXPOSE 8080

CMD ["./mocket"]
