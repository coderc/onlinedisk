FROM golang:1.19

WORKDIR /app

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]

# docker build -t onlinedisk-backend-image:v1 .
# docker build --platform linux/amd64 -t onlinedisk-backend-image-amd:v1 .
# docker run --name onlinedisk-backend-container-01 -p 8080:8080 -d onlinedisk-backend-image:v1