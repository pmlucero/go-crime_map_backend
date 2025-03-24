#!/bin/bash

# Cargar variables de entorno
source .env

# Crear la base de datos si no existe
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 || \
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -c "CREATE DATABASE $DB_NAME"

# Aplicar el esquema
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f internal/infrastructure/database/schema.sql 