-- Crear la tabla de ubicaciones
CREATE TABLE IF NOT EXISTS locations (
    id SERIAL PRIMARY KEY,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Eliminar tablas existentes si es necesario
DROP TABLE IF EXISTS crimes CASCADE;

-- Crear la tabla de delitos
CREATE TABLE crimes (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    crime_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Crear índices para mejorar el rendimiento de las consultas
CREATE INDEX IF NOT EXISTS idx_crimes_type ON crimes(type);
CREATE INDEX IF NOT EXISTS idx_crimes_date ON crimes(date);
CREATE INDEX IF NOT EXISTS idx_crimes_status ON crimes(status);
CREATE INDEX IF NOT EXISTS idx_locations_coordinates ON locations USING GIST (ll_to_earth(latitude, longitude));

-- Crear índices
CREATE INDEX idx_crimes_status ON crimes(status);
CREATE INDEX idx_crimes_crime_type ON crimes(crime_type);
CREATE INDEX idx_crimes_created_at ON crimes(created_at);
CREATE INDEX idx_crimes_deleted_at ON crimes(deleted_at);

-- Crear función para actualizar el timestamp updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Crear trigger para actualizar updated_at
CREATE TRIGGER update_crimes_updated_at
    BEFORE UPDATE ON crimes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_locations_updated_at
    BEFORE UPDATE ON locations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column(); 