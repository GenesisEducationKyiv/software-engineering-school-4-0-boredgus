FROM golang:1.22-alpine

WORKDIR /app

COPY . /app
RUN go build -o ./bin/daemon ./cmd/daemon/main.go
CMD [ "./bin/daemon" ]
