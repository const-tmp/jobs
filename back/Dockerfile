FROM golang:1.17

WORKDIR /jobs


COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build

CMD ./back
