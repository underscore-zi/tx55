FROM golang:latest AS builder
WORKDIR /go/src/app

COPY . .
RUN for dir in $(ls ./cmd); do go build -o /app/$dir tx55/cmd/$dir; done

FROM debian:stable-slim
RUN mkdir /app
COPY --from=builder /app /app
RUN apt-get update && apt-get install -y ca-certificates

HEALTHCHECK CMD /app/healthcheck
ENTRYPOINT ["/app/gameserver"]