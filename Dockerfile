# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY config.yml ./

RUN go build -o HellYeahPlayBot


CMD [ "/app/HellYeahPlayBot" ]