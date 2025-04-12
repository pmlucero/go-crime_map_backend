# Hoja de Ruta - Crime Map Backend

## Estado Actual
- ✅ Estructura base del proyecto con Clean Architecture
- ✅ Implementación inicial del caso de uso para crear delitos
- ✅ Configuración básica del servidor HTTP con Gin
- ✅ Endpoint de health check
- ✅ Endpoint básico para crear delitos
- ✅ Implementación de la generación de IDs usando UUID v4
- ✅ Implementación del repositorio PostgreSQL
- ✅ Implementación de endpoints CRUD completos
- ✅ Implementación de estadísticas de delitos
- ✅ Implementación de paginación y filtros básicos
- ✅ Pruebas de integración básicas
- ✅ Configuración de pre-commit hooks
- ✅ Implementación de CI/CD básico
- ✅ Documentación Swagger actualizada
- ✅ Validaciones básicas de datos
- ✅ Tests unitarios para casos de uso
- ✅ Tests de integración mejorados

## Pendientes Inmediatos

### 1. Validaciones de Negocio
- [x] Implementar validación de tipos de delitos permitidos
- [ ] Agregar validación de fechas (no futuras, no muy antiguas)
- [x] Implementar validación de ubicaciones (dentro de límites geográficos)
- [x] Agregar validación de descripciones (longitud mínima/máxima)
- [x] Implementar sanitización de datos de entrada

### 2. Pruebas
- [x] Implementar pruebas unitarias para:
  - [x] Casos de uso
  - [ ] Repositorios
  - [x] Controladores
- [x] Mejorar pruebas de integración
- [x] Configurar cobertura de código
- [ ] Implementar pruebas de carga

### 3. Documentación
- [x] Documentar la API con Swagger/OpenAPI
- [ ] Crear guía de instalación y configuración
- [ ] Documentar la estructura del proyecto
- [ ] Agregar ejemplos de uso
- [ ] Crear guía de contribución

### 4. Seguridad
- [ ] Implementar autenticación básica con API Key
  - [ ] Sistema de generación y revocación de API Keys
  - [ ] Almacenamiento seguro de API Keys
  - [ ] Middleware de validación de API Key
  - [ ] Documentación de uso de API Keys
- [ ] Implementar autorización por roles
  - [ ] Definir roles básicos (admin, user, read-only)
  - [ ] Implementar middleware de autorización
  - [ ] Documentar permisos por rol
- [ ] Agregar rate limiting por API Key
- [ ] Implementar CORS
- [ ] Implementar logging seguro
- [ ] (Backlog) Implementar OAuth2
  - [ ] Integración con proveedores OAuth2
  - [ ] Manejo de tokens JWT
  - [ ] Refresh tokens
  - [ ] Revocación de tokens
  - [ ] Migración de API Keys a OAuth2

### 5. Monitoreo y Logging
- [ ] Implementar sistema de logging estructurado
- [ ] Agregar métricas de la aplicación
- [ ] Configurar alertas
- [ ] Implementar tracing distribuido

## Próximos Pasos Sugeridos
1. Implementar validación de fechas
2. Completar pruebas unitarias para repositorios
3. Crear guía de instalación y configuración
4. Implementar autenticación con API Key
5. Configurar monitoreo y logging

## Consideraciones Técnicas
- Mantener la arquitectura limpia y los principios SOLID
- Asegurar la escalabilidad del sistema
- Implementar manejo de errores consistente
- Mantener la documentación actualizada
- Seguir las mejores prácticas de Go
- Implementar seguridad por capas (defense in depth)
- Considerar la migración futura a OAuth2 en el diseño actual 