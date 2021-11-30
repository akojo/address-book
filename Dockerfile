FROM golang:1.17-alpine AS build
WORKDIR /go/src/app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/address-book -ldflags='-w -s'

FROM gcr.io/distroless/static
COPY --from=build /go/bin/address-book / 
CMD ["/address-book"]
