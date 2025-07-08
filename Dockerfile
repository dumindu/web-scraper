FROM golang:1.24-alpine
WORKDIR /app

RUN apk --update-cache add gcc musl-dev tzdata

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/api ./cmd/api \
    && go build -ldflags '-w -s' -a -o ./bin/worker ./cmd/worker \
    && go build -ldflags '-w -s' -a -o ./bin/dbmigrate ./cmd/dbmigrate

EXPOSE 8080
CMD ["/app/bin/api"]
