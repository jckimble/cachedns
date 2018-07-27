FROM golang:latest as gobuild
ADD . /go/src/gitlab.com/paranoidsecurity/cachedns
WORKDIR /go/src/gitlab.com/paranoidsecurity/cachedns
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cachedns .

FROM scratch
WORKDIR /
COPY --from=gobuild /go/src/gitlab.com/paranoidsecurity/cachedns/cachedns .

EXPOSE 5353
ENTRYPOINT ["/cachedns"]
