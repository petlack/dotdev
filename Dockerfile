FROM golang:1.23.6-alpine3.21 AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o dotdev .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /app/dotdev /bin/dotdev
ENTRYPOINT ["dotdev"]
CMD ["dotdev"]
