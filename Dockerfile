FROM ubuntu:18.04 AS release

MAINTAINER gel0


ENV PGVER 10

# Обновление списка пакетов
RUN apt-get update
RUN apt-get install -y postgresql-$PGVER
RUN apt-get install -y curl gnupg2

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Этот параметр определяет, сколько памяти будет выделяться 
# postgres для кеширования данных (25% от всей оперативной памяти)
RUN echo "shared_buffers = 256MB" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Этот параметр помогает планировщику postgres определить количество 
# доступной памяти для дискового кеширования.
# На основе того, доступна память или нет,планировщик будет делать выбор между 
# использованием индексов и использованием сканирования таблицы. (75% от всей оперативной памяти)
RUN echo "effective_cache_size = 500MB" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Чекпойнт — это набор операций, которые выполняет postgres для гарантии того,
# что все изменения были записаны в файлы данных (следовательно при сбое, 
# восстановление происходит по последнему чекпойнту).
# RUN echo "checkpoint_=segments = 128MB" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Количество оперативной памяти, которое будет выделено на выполнение каждой операции
RUN echo "work_mem = 64MB" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Этот параметр определяет количество памяти для различных
# статистических и управляющих процессов (например вакуумизация). 
# - вообще хз что это, причем я отключаю автоматическую вакуумизацию
# при создании таблицы (autovacuum_enabled = false)
# RUN echo "maintainance_work_mem = 128MB" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Тупа для больших систем
RUN echo "wal_buffers = 1MB" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Отключение всего, что замедляет работу
RUN echo "synchronous_commit = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "fsync = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "full_page_writes = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "maintenance_work_mem = 300MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_statement = 'none'" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Остальное под железо машины
RUN echo "max_connections = 200" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "checkpoint_completion_target = 0.9" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "default_statistics_target = 100" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "random_page_cost = 1.1" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "effective_io_concurrency = 200" >> /etc/postgresql/$PGVER/main/postgresql.conf

RUN echo "min_wal_size = 512MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "max_wal_size = 1GB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "max_worker_processes = 8" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "max_parallel_workers_per_gather = 4" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "max_parallel_workers = 4" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "unix_socket_directories = '/var/run/postgresql'" >> /etc/postgresql/$PGVER/main/postgresql.conf


EXPOSE 5432

# RUN apt-get update && apt-get install -y postgresql-$PGVER

RUN apt-get install -y wget

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql -c "CREATE USER docker WITH SUPERUSER PASSWORD '1337';" &&\
    createdb -O docker forum &&\
    psql -c "GRANT ALL ON DATABASE forum TO docker;" &&\
    psql -d forum -c "CREATE EXTENSION IF NOT EXISTS citext;" &&\
    /etc/init.d/postgresql stop

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root

RUN wget https://storage.googleapis.com/golang/go1.11.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.11.linux-amd64.tar.gz
RUN apt-get install -y git


# Выставляем переменную окружения для сборки проекта
ENV GOPATH /opt/go

ENV PATH $PATH:/usr/local/go/bin

RUN go get github.com/jackc/pgx
RUN go get github.com/gin-gonic/gin

ADD / /

EXPOSE 5000

CMD service postgresql start && go run main.go

