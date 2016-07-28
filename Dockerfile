FROM ubuntu

RUN mkdir /app
ADD dist/server /app/server

EXPOSE 8085

CMD ["/app/server"]