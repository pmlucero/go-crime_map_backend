#!/bin/bash

# Ejecutar tests y generar reporte de cobertura para todos los paquetes
echo "Ejecutando tests y generando reporte de cobertura..."
go test -coverpkg=./... ./... -coverprofile=coverage.out && go tool cover -func=coverage.out

# Limpiar archivo temporal
rm -f coverage.out 