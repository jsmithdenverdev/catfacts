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

RUN go build -o bin/sendfact ./cmd/sendfact/*.go

### App
FROM alpine as app
COPY --from=base app /app

ENV DATABASE_URL=""
ENV TWILIO_SID=""
ENV TWILIO_TOKEN=""
ENV TWILIO_FROM=""

ENTRYPOINT ["/app/bin/sendfact"]