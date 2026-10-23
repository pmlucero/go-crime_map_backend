-- Primero hacemos backup de los datos existentes
CREATE TABLE crimes_backup AS SELECT * FROM crimes;

-- Eliminamos la tabla existente
DROP TABLE crimes;

-- Creamos la nueva tabla con la estructura actualizada
CREATE TABLE crimes (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    crime_type VARCHAR(50) NOT NULL,
    description TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    status VARCHAR(20) NOT NULL,
    address TEXT NOT NULL,
    address_number VARCHAR(50),
    city VARCHAR(100),
    province VARCHAR(100),
    country VARCHAR(100),
    zip_code VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Creamos los índices necesarios
CREATE INDEX idx_crimes_uuid ON crimes(uuid);
CREATE INDEX idx_crimes_status ON crimes(status);
CREATE INDEX idx_crimes_created_at ON crimes(created_at);
CREATE INDEX idx_crimes_location ON crimes USING gist (
    ll_to_earth(latitude, longitude)
);

-- Restauramos los datos con la nueva estructura (asignando uuid nuevos)
INSERT INTO crimes (
    id, uuid, crime_type, description, latitude, longitude,
    status, address, address_number, city, province, country, zip_code,
    created_at, updated_at, deleted_at
)
SELECT 
    id,
    gen_random_uuid()::text,
    COALESCE(crime_type, 'UNKNOWN'),
    description,
    COALESCE(latitude, 0),
    COALESCE(longitude, 0),
    COALESCE(status, 'ACTIVE'),
    address,
    address_number,
    city,
    province,
    country,
    zip_code,
    created_at,
    updated_at,
    deleted_at
FROM crimes_backup;

-- Eliminamos la tabla de backup
DROP TABLE crimes_backup; 