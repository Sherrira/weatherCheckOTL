FROM golang:1.21 as build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun ./cmd

FROM scratch
WORKDIR /app
COPY ./cmd/.env .
COPY --from=build /app/cloudrun .
ENTRYPOINT ["./cloudrun"]