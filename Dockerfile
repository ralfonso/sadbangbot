FROM alpine:3.4
RUN apk --update upgrade && \
    apk add curl ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*
COPY deploy/sadbangbotd /sadbangbotd
ENTRYPOINT ["/sadbangbotd"]
