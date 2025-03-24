#!/bin/bash

# Inicializar la base de datos de producción
echo "Inicializando base de datos de producción..."
psql -U $USER postgres -f scripts/init_db.sql

echo "Base de datos inicializada correctamente" 