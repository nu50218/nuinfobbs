FROM golang:1.13 as build

RUN go get "github.com/dghubble/go-twitter/twitter"
RUN go get "github.com/dghubble/oauth1"
RUN go get "cloud.google.com/go/firestore"
RUN go get "firebase.google.com/go"
RUN go get "google.golang.org/api/iterator"
RUN go get "github.com/caarlos0/env"
RUN go get "github.com/nu50218/nuinfobbs/library/jobutils"
WORKDIR /go/src/github.com/nu50218/nuinfobbs/containers/twitter/src
COPY src .
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go tweet.go

FROM alpine:3.10
RUN apk add --no-cache ca-certificates
COPY --from=build /go/src/github.com/nu50218/nuinfobbs/containers/twitter/src/app app

CMD [ "./app" ]
