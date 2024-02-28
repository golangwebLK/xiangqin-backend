FROM alpine

WORKDIR /app

COPY xiangqin-backend /usr/local/bin/

EXPOSE 8000

RUN chmod +x /usr/local/bin/xiangqin-backend

CMD ["/usr/local/bin/xiangqin-backend","server"]
