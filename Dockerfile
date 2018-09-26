FROM golang:1.11-alpine3.8 as build

WORKDIR /go/src/github.com/vmj/redirect

COPY redirect.go ./
RUN CGO_ENABLED=0 go build -a -o redirect

FROM scratch
COPY --from=build /go/src/github.com/vmj/redirect/redirect /
CMD ["/redirect", "-protocol", "https", "-port", ""]
