FROM golang:1.18-alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY ./src ./src

RUN go build -o ./bin/hello-world ./src

EXPOSE 8080

CMD [ "./bin/hello-world" ]