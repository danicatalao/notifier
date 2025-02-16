FROM golang:1.23-alpine AS dependencies
COPY go.mod go.sum /dependencies/
WORKDIR /dependencies
RUN go mod download

FROM dependencies AS build
COPY --from=dependencies /go/pkg /go/pkg
COPY . /server
WORKDIR /server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/consumer ./cmd/consumer
RUN apk --no-cache add ca-certificates tzdata

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/consumer /consumer
CMD ["/consumer"]