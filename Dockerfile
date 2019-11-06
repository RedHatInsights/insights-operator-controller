FROM golang:1.13 AS builder

ARG SOURCE_REPOSITORY_URL

RUN git clone $SOURCE_REPOSITORY_URL && \
    cd insights-operator-controller && \
    go build 

FROM registry.access.redhat.com/ubi8-minimal

COPY --from=builder /go/insights-operator-controller/insights-operator-controller .

RUN curl -L https://github.com/Yelp/dumb-init/releases/download/v1.2.2/dumb-init_1.2.2_amd64 -o /usr/bin/dumb-init && \
    chmod a+x /usr/bin/dumb-init && \
    chmod a+x /insights-operator-controller

USER 1001

CMD ["dumb-init", "/insights-operator-controller"]
