# Hoja de Ruta - Crime Map Backend

## Estado Actual
- ✅ Estructura base del proyecto con Clean Architecture
- ✅ Implementación inicial del caso de uso para crear delitos
- ✅ Configuración básica del servidor HTTP con Gin
- ✅ Endpoint de health check
- ✅ Endpoint básico para crear delitos

## Pendientes Inmediatos

### 1. Implementación de la Generación de IDs
- [ ] Implementar generación de IDs únicos usando UUID v4
- [ ] Agregar validación de formato de ID
- [ ] Considerar implementar un servicio de generación de IDs distribuido

### 2. Implementación del Repositorio
- [ ] Definir la base de datos a utilizar (PostgreSQL recomendado)
- [ ] Crear la estructura de la base de datos
- [ ] Implementar el repositorio concreto para PostgreSQL
- [ ] Agregar migraciones de base de datos
- [ ] Implementar manejo de conexiones y pool de conexiones

### 3. Validaciones de Negocio
- [ ] Implementar validación de tipos de delitos permitidos
- [ ] Agregar validación de fechas (no futuras, no muy antiguas)
- [ ] Implementar validación de ubicaciones (dentro de límites geográficos)
- [ ] Agregar validación de descripciones (longitud mínima/máxima)
- [ ] Implementar sanitización de datos de entrada

### 4. Endpoints Adicionales
- [ ] Implementar GET /api/v1/crimes (listar todos los delitos)
- [ ] Implementar GET /api/v1/crimes/:id (obtener un delito específico)
- [ ] Implementar PUT /api/v1/crimes/:id (actualizar un delito)
- [ ] Implementar DELETE /api/v1/crimes/:id (eliminar un delito)
- [ ] Agregar paginación para el listado de delitos
- [ ] Implementar filtros por tipo, fecha y ubicación

### 5. Pruebas
- [ ] Configurar el entorno de pruebas
- [ ] Implementar pruebas unitarias para:
  - [ ] Casos de uso
  - [ ] Repositorios
  - [ ] Controladores
- [ ] Implementar pruebas de integración
- [ ] Configurar cobertura de código
- [ ] Implementar pruebas de carga

### 6. Documentación
- [ ] Documentar la API con Swagger/OpenAPI
- [ ] Crear guía de instalación y configuración
- [ ] Documentar la estructura del proyecto
- [ ] Agregar ejemplos de uso
- [ ] Crear guía de contribución

### 7. Seguridad
- [ ] Implementar autenticación
- [ ] Implementar autorización
- [ ] Agregar rate limiting
- [ ] Implementar CORS
- [ ] Agregar validación de tokens JWT
- [ ] Implementar logging seguro

### 8. Monitoreo y Logging
- [ ] Implementar sistema de logging estructurado
- [ ] Agregar métricas de la aplicación
- [ ] Configurar alertas
- [ ] Implementar tracing distribuido

## Próximos Pasos Sugeridos
1. Implementar la generación de IDs únicos
2. Configurar y probar la base de datos
3. Implementar el repositorio concreto
4. Agregar validaciones de negocio
5. Implementar los endpoints faltantes
6. Configurar el entorno de pruebas

## Consideraciones Técnicas
- Mantener la arquitectura limpia y los principios SOLID
- Asegurar la escalabilidad del sistema
- Implementar manejo de errores consistente
- Mantener la documentación actualizada
- Seguir las mejores prácticas de Go 