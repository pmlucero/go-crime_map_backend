#!/bin/bash

# Variables de entorno para la base de datos de test
export TEST_DB_HOST=${TEST_DB_HOST:-localhost}
export TEST_DB_PORT=${TEST_DB_PORT:-5432}
export TEST_DB_USER=${TEST_DB_USER:-postgres}
export TEST_DB_PASSWORD=${TEST_DB_PASSWORD:-postgres}
export TEST_DB_NAME=${TEST_DB_NAME:-crime_map_test}

# Crear la base de datos de test
PGPASSWORD=$TEST_DB_PASSWORD psql -h $TEST_DB_HOST -p $TEST_DB_PORT -U $TEST_DB_USER -d postgres -c "DROP DATABASE IF EXISTS $TEST_DB_NAME;"
PGPASSWORD=$TEST_DB_PASSWORD psql -h $TEST_DB_HOST -p $TEST_DB_PORT -U $TEST_DB_USER -d postgres -c "CREATE DATABASE $TEST_DB_NAME;"

# Habilitar la extensión earthdistance
PGPASSWORD=$TEST_DB_PASSWORD psql -h $TEST_DB_HOST -p $TEST_DB_PORT -U $TEST_DB_USER -d $TEST_DB_NAME -c "CREATE EXTENSION IF NOT EXISTS cube;"
PGPASSWORD=$TEST_DB_PASSWORD psql -h $TEST_DB_HOST -p $TEST_DB_PORT -U $TEST_DB_USER -d $TEST_DB_NAME -c "CREATE EXTENSION IF NOT EXISTS earthdistance;"

# Ejecutar el esquema en la base de datos de test
PGPASSWORD=$TEST_DB_PASSWORD psql -h $TEST_DB_HOST -p $TEST_DB_PORT -U $TEST_DB_USER -d $TEST_DB_NAME -f internal/infrastructure/database/schema.sql

echo "Base de datos de test configurada exitosamente" 