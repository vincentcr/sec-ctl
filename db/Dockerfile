FROM postgres:9.5-alpine

ARG DB_PASSWORD
ENV DB_NAME tpimon_dev
ENV DB_USER tpimon_dev

RUN  mkdir /pgsql-bootstrap

# save the hashed password to a file. it will be picked later by the the init-db.sh script.
RUN echo ${DB_PASSWORD}${DB_USER} | md5sum | cut -d' ' -f1 > /pgsql-bootstrap/md5_passwd

COPY init-db.sh /docker-entrypoint-initdb.d/
COPY *.sql /pgsql-bootstrap/


