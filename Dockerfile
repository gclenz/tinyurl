FROM golang:1.24 AS build
WORKDIR /app
COPY . /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o api

FROM alpine
RUN apk add bash
RUN apk add curl
WORKDIR /app
COPY --from=build /app/api ./
EXPOSE 8080
CMD [ "./api" ]