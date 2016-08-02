FROM centurylink/ca-certs

ADD ./dist/server /server

EXPOSE 8085

CMD ["/server"]