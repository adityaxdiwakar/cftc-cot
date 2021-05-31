From golang:1.15-buster as builder

WORKDIR /go/src/cftc-cot/api
COPY . /go/src/cftc-cot

RUN go get -d -v
RUN go build -o /go/bin/cftc-cot

FROM gcr.io/distroless/base-debian10
COPY --from=builder /go/bin/cftc-cot /
CMD ["/cftc-cot"]
