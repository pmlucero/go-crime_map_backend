#!/bin/bash

# Inicializar la base de datos de test
echo "Inicializando base de datos de test..."
psql -U $USER postgres -f scripts/init_test_db.sql

echo "Base de datos de test inicializada correctamente" 