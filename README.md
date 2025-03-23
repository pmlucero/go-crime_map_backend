# Crime Map Backend

Backend para la aplicación Crime Map, implementado en Go siguiendo los principios de Clean Architecture.

## Estructura del Proyecto

```
.
├── cmd/
│   └── api/              # Punto de entrada de la aplicación
├── internal/
│   ├── domain/          # Entidades y reglas de negocio
│   │   ├── entities/    # Entidades del dominio
│   │   └── repositories/# Interfaces de repositorios
│   ├── usecases/        # Casos de uso de la aplicación
│   ├── interfaces/      # Adaptadores y controladores
│   │   ├── http/       # Controladores HTTP
│   │   └── repositories/# Implementaciones de repositorios
│   └── infrastructure/  # Implementaciones técnicas
│       ├── config/     # Configuración
│       ├── database/   # Conexión a base de datos
│       └── server/     # Servidor HTTP
└── pkg/
    └── utils/          # Utilidades compartidas
```

## Requisitos

- Go 1.21 o superior
- Docker (opcional, para desarrollo)

## Instalación

1. Clonar el repositorio:
```bash
git clone https://github.com/tu-usuario/go-crime_map_backend.git
cd go-crime_map_backend
```

2. Instalar dependencias:
```bash
go mod download
```

3. Ejecutar la aplicación:
```bash
go run cmd/api/main.go
```

## Desarrollo

Para ejecutar el servidor en modo desarrollo:
```bash
go run cmd/api/main.go
```

El servidor estará disponible en `http://localhost:8080`

## Endpoints

- `GET /health`: Verificar el estado del servidor

## Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.
