FROM golang:alpine AS builder
LABEL authors = "Aidana_BK & Erlangt & LevapMik"
LABEL project-name = "forum"
WORKDIR /app
COPY . .
RUN apk add build-base && go build -o forum cmd/main.go
FROM alpine
WORKDIR /app 
COPY --from=builder /app .
CMD ["./forum"]
EXPOSE 8081