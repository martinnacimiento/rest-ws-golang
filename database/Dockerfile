FROM postgres:12.1

COPY ./init.sql /docker-entrypoint-initdb.d/1.sql

CMD ["postgres"]