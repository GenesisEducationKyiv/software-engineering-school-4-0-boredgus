# STEP 1: build
FROM golang:1.22-alpine as build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . /app/

RUN go build -o /bin/server ./cmd/server/main.go


# STEP 2: run binary executable
FROM scratch

WORKDIR /

COPY --from=build /bin/server /server

ENTRYPOINT [ "/server" ]
