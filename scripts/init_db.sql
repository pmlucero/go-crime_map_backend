-- Crear la base de datos si no existe
CREATE DATABASE crime_map;

-- Conectar a la base de datos
\c crime_map

-- Crear la extensión uuid-ossp si no existe
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Leer y ejecutar el esquema
\i internal/infrastructure/database/schema.sql 