FROM golang:1.15 AS build
COPY . /app
WORKDIR /app
RUN go get -d
RUN go build -o magazi

EXPOSE 9242

ENTRYPOINT [ "./magazi"]
