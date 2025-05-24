-- Hacemos backup de los datos existentes
CREATE TABLE crimes_backup AS SELECT * FROM crimes;

-- Eliminamos la tabla existente
DROP TABLE crimes;

-- Recreamos la tabla con la estructura anterior
CREATE TABLE crimes (
    id BIGSERIAL PRIMARY KEY,
    crime_type VARCHAR(50),
    description TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    title VARCHAR(100),
    location VARCHAR(200),
    status VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Restauramos los datos con la estructura anterior
INSERT INTO crimes (
    id, crime_type, description, latitude, longitude,
    status, created_at, updated_at, deleted_at
)
SELECT 
    id, type, description, latitude, longitude,
    status, created_at, updated_at, deleted_at
FROM crimes_backup;

-- Eliminamos la tabla de backup
DROP TABLE crimes_backup; 