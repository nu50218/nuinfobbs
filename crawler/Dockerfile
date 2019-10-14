FROM golang:1.13 as build

WORKDIR /go/src/github.com/nu50218/nuinfobbs/crawler/src
COPY src .
RUN go get "github.com/go-sql-driver/mysql"
RUN go get "github.com/PuerkitoBio/goquery"
RUN CGO_ENABLED=0 GOOS=linux go build -o crawler main.go crawl.go

FROM alpine:3.10
RUN apk add --no-cache ca-certificates
COPY --from=build /go/src/github.com/nu50218/nuinfobbs/crawler/src/crawler crawler
