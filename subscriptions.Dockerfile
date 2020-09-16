FROM golang as base

WORKDIR /app

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o bin/subscriptions ./cmd/subscriptions/*.go

### App
FROM alpine as app
COPY --from=base app /app

ENV DATABASE_URL=""
EXPOSE 8080

ENTRYPOINT ["/app/bin/subscriptions"]
