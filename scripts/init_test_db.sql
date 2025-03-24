-- Crear la base de datos si no existe
CREATE DATABASE crime_map_test;

-- Conectar a la base de datos
\c crime_map_test;

-- Crear el esquema test si no existe
CREATE SCHEMA IF NOT EXISTS test;

-- Establecer el esquema por defecto
SET search_path TO test;

-- Crear la extensión uuid-ossp si no existe
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Leer y ejecutar el esquema
\i internal/infrastructure/database/schema.sql 