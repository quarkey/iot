FROM mdillon/postgis:latest

ENV TERM linux
ENV DEBIAN_FRONTEND noninteractive

WORKDIR /root

RUN echo "host all  all    0.0.0.0/0  md5" >> /var/lib/postgresql/data/pg_hba.conf
RUN echo "listen_addresses='*'" >> /var/lib/postgresql/datapostgresql.conf

EXPOSE 5432
