FROM golang:alpine as builder

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o oidc-thumbprint .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/oidc-thumbprint .

ENTRYPOINT ["/oidc-thumbprint"]
