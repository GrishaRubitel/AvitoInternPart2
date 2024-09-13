FROM golang:1.22.2-alpine AS builder
WORKDIR /usr/local/src
RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["go/go.mod", "go/go.sum", "./"]
RUN go mod download

COPY go ./
RUN go build -o ./bin/go main.go tenderApiHandler.go bidApiHandler.go universalFuncs.go


FROM alpine
COPY .env .env
COPY --from=builder /usr/local/src/bin/go /
EXPOSE 8080
CMD [ "/go" ]