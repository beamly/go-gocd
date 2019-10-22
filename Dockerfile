ARG GOCD_VERSION
FROM gocd/gocd-server:${GOCD_VERSION}

ARG UID

USER root

RUN apk --no-cache add shadow && \
    usermod -u ${UID} go

USER go
