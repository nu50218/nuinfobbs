FROM golang:1.13 as build

WORKDIR /go/src/github.com/nu50218/nuinfobbs/containers/twitter/src
RUN go get "github.com/dghubble/go-twitter/twitter"
RUN go get "github.com/dghubble/oauth1"
RUN go get "cloud.google.com/go/firestore"
RUN go get "firebase.google.com/go"
RUN go get "google.golang.org/api/iterator"
COPY src .
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

FROM alpine:3.10
RUN apk add --no-cache ca-certificates
COPY --from=build /go/src/github.com/nu50218/nuinfobbs/containers/twitter/src/app app

CMD [ "./app" ]