FROM centurylink/ca-certs

RUN mkdir /app
ADD dist/server /app/server

EXPOSE 8085

CMD ["/app/server"]