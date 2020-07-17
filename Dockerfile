FROM alpine:latest AS builder
ARG VERSION
ADD /binaries/monitoror-linux-amd64-${VERSION} /bin/monitoror
RUN chmod +x /bin/monitoror

FROM alpine:latest
RUN apk update && \
    apk --no-cache add ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*
COPY --from=builder /bin/monitoror /bin/monitoror
EXPOSE 8080
CMD [ "/bin/monitoror" ]
