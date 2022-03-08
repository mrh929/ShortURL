FROM golang:latest

ENV WORKDIR '/src'
COPY ./entrypoint.sh /
RUN mkdir -p ${WORKDIR} && chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]