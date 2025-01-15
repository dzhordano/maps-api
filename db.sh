#!/bin/bash

psql -v ON_ERROR_STOP=1 --username "${POSTGRES_USER}" --dbname "${POSTGRES_DB}" <<-EOSQL
    CREATE EXTENSION postgis;
    ALTER DATABASE postgres REFRESH COLLATION VERSION;
EOSQL
