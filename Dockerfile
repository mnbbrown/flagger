FROM alpine

FROM scratch
COPY --from=0 /etc/passwd /etc/passwd
ADD flagctl-linux-amd64 /flagctl
ADD ui /ui

USER 405:405
EXPOSE 8082
ENTRYPOINT ["/flagctl"]
CMD ["serve"]
