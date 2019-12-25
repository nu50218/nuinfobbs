FROM golang:1.13 as build

WORKDIR /src
ADD src/go.mod src/go.sum /src/
RUN go mod download

COPY src .
RUN CGO_ENABLED=0 GOOS=linux go build -o app *.go

FROM alpine:3.10
RUN apk add --no-cache ca-certificates
COPY --from=build /src app

CMD [ "./app" ]
