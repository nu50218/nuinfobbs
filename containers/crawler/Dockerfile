FROM golang:1.13 as build

WORKDIR /src
ADD src/go.mod src/go.sum /src/
RUN go mod download

COPY src .
RUN CGO_ENABLED=0 GOOS=linux go build -o app *.go

FROM nu50218/cloud-run-command:3.10

# https://qiita.com/_mkazutaka/items/12f8d560e163b8ab1813
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /src/app app

CMD [ "cloud-run-command","./app" ]
