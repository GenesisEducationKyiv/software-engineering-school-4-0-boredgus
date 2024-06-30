FROM golang:1.22-alpine

WORKDIR /app

COPY . /app
RUN go build -o ./bin/server ./cmd/server/main.go
CMD [ "./bin/server" ]
