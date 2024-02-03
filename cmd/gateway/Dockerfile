FROM golang:latest as builder

WORKDIR /aurora

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM apline:latesr

COPY --from=builder /aurora/main /
ENTRYPOINT ["main"]