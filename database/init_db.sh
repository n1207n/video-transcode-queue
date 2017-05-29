#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE USER video_app;
    CREATE DATABASE video_app;
    GRANT ALL PRIVILEGES ON DATABASE video_app TO video_app;
EOSQL
