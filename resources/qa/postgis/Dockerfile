FROM postgis
WORKDIR /root
COPY init-user-db.sh /docker-entrypoint-initdb.d/
EXPOSE 5432

docker run -d --name some-other-postgres -e POSTGRES_PASSWORD=mypassword -p 25432:5432 postgis/postgis