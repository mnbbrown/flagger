FROM alpine
USER nobody:nobody
ADD flagger /flagger
ENTRYPOINT /flagger
