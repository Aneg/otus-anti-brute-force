FROM golang:1.14

WORKDIR /app
#COPY . .
COPY go.mod go.sum ./
RUN go mod download

EXPOSE 8080