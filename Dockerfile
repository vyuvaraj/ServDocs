FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o servdocs

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/servdocs .
COPY --from=builder /app/docs.html .
EXPOSE 3000
ENTRYPOINT ["./servdocs"]
