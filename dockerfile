FROM alpine:3.18

WORKDIR /app

COPY go-bin ./go-bin
RUN chmod +x ./go-bin

COPY env.json ./env.json
COPY frontend ./frontend

COPY localhost.pem ./localhost.pem
COPY localhost-key.pem ./localhost-key.pem

CMD ["/app/go-bin"]

EXPOSE 80