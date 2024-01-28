FROM golang:1.21.6-alpine3.19 as build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o /app image-service/cmd/imageservice

FROM alpine:3.19
COPY --from=build /app /app
EXPOSE 1111
CMD ["/app"]