FROM golang:1.16.3-alpine

WORKDIR /app

COPY . .

EXPOSE 3000

CMD [ "go", "run", "cmd/web/main.go" ]