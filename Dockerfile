FROM golang:alpine3.17 as build
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o build/ordersys cmd/ordersys/main.go

FROM scratch
WORKDIR /app
COPY --from=build /app/build/ordersys .
CMD [ "./ordersys" ]
