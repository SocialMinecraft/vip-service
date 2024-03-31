FROM golang:1.22

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./

CMD ["app"]

#docker build --platform linux/amd64 --no-cache -t dmgarvis/vip-service:latest .
#docker push dmgarvis/vip-service:latest

#docker run --env-file .env dmgarvis/vip-service:latest