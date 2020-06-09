FROM alpine:3.9.6 as dev
WORKDIR /usr/src/app
COPY main /usr/local/bin/faas
EXPOSE 8081
CMD /usr/local/bin/faas
