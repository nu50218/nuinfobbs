FROM golang:1.13 as build

WORKDIR /go/src/github.com/nu50218/nuinfobbs/app/src
COPY src .
RUN go get "github.com/go-sql-driver/mysql"
RUN go get "github.com/dghubble/go-twitter/twitter"
RUN go get "github.com/dghubble/oauth1"
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go tweet.go

FROM alpine:3.10
RUN apk add --no-cache ca-certificates
COPY --from=build /go/src/github.com/nu50218/nuinfobbs/app/src/app app
