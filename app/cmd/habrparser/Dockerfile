FROM golang:1.16

WORKDIR /usr/src/habrparser

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
WORKDIR /usr/src/habrparser/cmd/habrparser
RUN go build -v -o /usr/local/bin/habrparser

CMD ["habrparser"]
