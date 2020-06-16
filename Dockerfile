FROM alpine:latest AS builder
ARG VERSION
ADD /binaries/monitoror-linux-amd64-${VERSION} /bin/monitoror
RUN chmod +x /bin/monitoror

FROM scratch
COPY --from=builder /bin/monitoror /bin/monitoror
EXPOSE 8080
CMD [ "/bin/monitoror" ]
