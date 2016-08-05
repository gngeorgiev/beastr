FROM centurylink/ca-certs

ADD ./dist/server /server

ENV REDIS_ADDRESS=redis:6379

EXPOSE 8085

CMD ["/server"]