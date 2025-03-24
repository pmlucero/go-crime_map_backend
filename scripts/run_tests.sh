#!/bin/bash

# Inicializar la base de datos de pruebas
echo "Inicializando base de datos de pruebas..."
psql -U $USER postgres -f scripts/init_test_db.sql

# Ejecutar los tests
echo "Ejecutando tests..."
go test -v ./... 