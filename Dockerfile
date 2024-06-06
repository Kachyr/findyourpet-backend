FROM golang:1.21.6 AS builder

WORKDIR /app

COPY go.mod go.sum Makefile ./

RUN go mod download

COPY . .

RUN make build-go
RUN chmod +x findyourpet-backend

FROM alpine:3.14 
# search c-alpine 
ENV GIN_PORT="3000"
ENV GIN_MODE="release"
ENV LOG_LEVEL=0
ENV ENV="PROD" 
ENV AWS_REGION="eu-north-1"


ENV SECRET=FINDYOURPET

RUN apk add libc6-compat

WORKDIR /executable
COPY --from=builder /app/findyourpet-backend .

EXPOSE 3000

CMD ["./findyourpet-backend"]