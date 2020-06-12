FROM alpine:3.12.0

COPY scripts /scripts

ENTRYPOINT [ "/scripts/entrypoint.sh" ]