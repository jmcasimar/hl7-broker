FROM alpine:latest

WORKDIR /app

COPY public .

COPY health24-gateway .

EXPOSE 8080
EXPOSE 4600

CMD ["./health24-gateway"]
