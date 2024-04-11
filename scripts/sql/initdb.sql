-- Ensure the PostGIS extension is installed
CREATE EXTENSION IF NOT EXISTS postgis;

-- Drop existing tables if they exist
DROP TABLE IF EXISTS soil_data CASCADE;
DROP TABLE IF EXISTS satellite_imagery CASCADE;
DROP TABLE IF EXISTS images CASCADE;
DROP TABLE IF EXISTS vineyards CASCADE;

-- Create vineyards table with a bounding box column
CREATE TABLE vineyards (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    bbox GEOMETRY(POLYGON, 4326),  -- POLYGON type with SRID 4326 (WGS 84)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create images table storing references to images stored in cloud storage
CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    vineyard_id INTEGER NOT NULL,
    image_url TEXT NOT NULL,
    bbox GEOMETRY(POLYGON, 4326),  -- Bounding box defining the extents of the image
    captured_at TIMESTAMP WITH TIME ZONE NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vineyard_id) REFERENCES vineyards (id) ON DELETE CASCADE
);

-- Create satellite imagery table storing references to satellite images stored in cloud storage
CREATE TABLE satellite_imagery (
    id SERIAL PRIMARY KEY,
    vineyard_id INTEGER NOT NULL,
    image_url TEXT NOT NULL,
    bbox GEOMETRY(POLYGON, 4326),  -- Bounding box for the satellite image
    resolution DECIMAL(10,2) DEFAULT 0.00,
    captured_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vineyard_id) REFERENCES vineyards (id) ON DELETE CASCADE
);

-- Create soil data table
CREATE TABLE soil_data (
    id SERIAL PRIMARY KEY,
    vineyard_id INTEGER NOT NULL,
    data JSONB,  -- JSONB for flexible data storage
    location GEOMETRY(POINT, 4326),  -- Exact point where soil data was collected
    sampled_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vineyard_id) REFERENCES vineyards (id) ON DELETE CASCADE
);
