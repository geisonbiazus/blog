FROM golang:1.16.3-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o blog cmd/web/main.go

FROM alpine:3

WORKDIR /app

COPY --from=builder /app/blog .
COPY --from=builder /app/web ./web
COPY --from=builder /app/posts ./posts

EXPOSE 3000

CMD [ "./blog" ]