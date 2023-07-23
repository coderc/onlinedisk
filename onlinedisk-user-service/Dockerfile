FROM golang:1.19

WORKDIR /app

COPY . .

RUN go build -o main .

EXPOSE 8081

CMD ["./main"]

# docker build -t onlinedisk-user-service-image:v1 .
# docker run --name onlinedisk-user-service-container-01 -p 8081:8081 -d onlinedisk-user-service-image:v1