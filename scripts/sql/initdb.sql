-- Ensure the PostGIS extension is installed
CREATE EXTENSION IF NOT EXISTS postgis;

-- Drop existing tables if they exist to prevent errors in table creation
DROP TABLE IF EXISTS weather_data CASCADE;
DROP TABLE IF EXISTS pest_data CASCADE;
DROP TABLE IF EXISTS soil_data CASCADE;
DROP TABLE IF EXISTS satellite_imagery CASCADE;
DROP TABLE IF EXISTS images CASCADE;
DROP TABLE IF EXISTS vineyards CASCADE;

-- Create vineyards table with a bounding box column
CREATE TABLE vineyards (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    bbox GEOMETRY(POLYGON, 4326),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create images table storing references to images stored in cloud storage
CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    vineyard_id INTEGER NOT NULL,
    image_url TEXT NOT NULL,
    description TEXT,
    bbox GEOMETRY(POLYGON, 4326),
    captured_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vineyard_id) REFERENCES vineyards(id) ON DELETE CASCADE
);

-- Create satellite imagery table storing references to satellite images stored in cloud storage
CREATE TABLE satellite_imagery (
    id SERIAL PRIMARY KEY,
    vineyard_id INTEGER NOT NULL,
    image_url TEXT NOT NULL,
    bbox GEOMETRY(POLYGON, 4326),
    resolution DECIMAL(10,2) DEFAULT 0.00,
    captured_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vineyard_id) REFERENCES vineyards(id) ON DELETE CASCADE
);

-- Create soil data table
CREATE TABLE soil_data (
    id SERIAL PRIMARY KEY,
    vineyard_id INTEGER NOT NULL,
    data JSONB NOT NULL,
    location GEOMETRY(POINT, 4326),
    sampled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vineyard_id) REFERENCES vineyards(id) ON DELETE CASCADE
);

-- Create Pest Data Table
CREATE TABLE pest_data (
    id SERIAL PRIMARY KEY,
    vineyard_id INTEGER NOT NULL,
    description TEXT,
    observation_date TIMESTAMP WITH TIME ZONE NOT NULL,
    location GEOMETRY(POINT, 4326),
    FOREIGN KEY (vineyard_id) REFERENCES vineyards(id) ON DELETE CASCADE
);

-- Create Weather Data Table
CREATE TABLE weather_data (
    id SERIAL PRIMARY KEY,
    vineyard_id INTEGER NOT NULL,
    temperature DECIMAL(4, 2) NOT NULL,
    humidity DECIMAL(4, 2) NOT NULL,
    observation_time TIMESTAMP WITH TIME ZONE NOT NULL,
    location GEOMETRY(POINT, 4326),
    FOREIGN KEY (vineyard_id) REFERENCES vineyards(id) ON DELETE CASCADE
);
