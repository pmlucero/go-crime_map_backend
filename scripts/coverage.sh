#!/bin/bash

# Crear directorio para reportes si no existe
mkdir -p coverage

# Ejecutar tests con cobertura
go test -coverprofile=coverage/coverage.out ./...

# Generar reporte HTML
go tool cover -html=coverage/coverage.out -o coverage/coverage.html

# Generar reporte en formato cobertura (para CI/CD)
gocover-cobertura < coverage/coverage.out > coverage/coverage.xml

# Mostrar resumen de cobertura
echo "Resumen de cobertura:"
go tool cover -func=coverage/coverage.out 