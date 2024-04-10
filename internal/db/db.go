package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"viticulture-harvester-app/internal/model"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Set up connection pool settings (if needed based on your application load)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error verifying connection to database: %w", err)
	}

	return &DB{db}, nil
}

// SaveVineyard inserts a new Vineyard record into the database.
func (db *DB) SaveVineyard(ctx context.Context, vineyard *model.Vineyard) error {
	// Assuming a simplified structure; adjust according to our schema
	const query = `
    INSERT INTO vineyards (name, location) 
    VALUES ($1, $2) 
    RETURNING id`
	err := db.QueryRowContext(ctx, query, vineyard.Name, vineyard.Location).Scan(&vineyard.ID)
	if err != nil {
		return fmt.Errorf("inserting vineyard: %w", err)
	}
	return nil
}

// GetVineyard retrieves a Vineyard by ID.
func (db *DB) GetVineyard(ctx context.Context, id int) (*model.Vineyard, error) {
	const query = `
    SELECT id, name, location 
    FROM vineyards 
    WHERE id = $1`
	vineyard := &model.Vineyard{}
	err := db.QueryRowContext(ctx, query, id).Scan(&vineyard.ID, &vineyard.Name, &vineyard.Location)
	if err != nil {
		return nil, fmt.Errorf("retrieving vineyard by ID: %w", err)
	}
	return vineyard, nil
}

// UpdateVineyard updates a given Vineyard's details.
func (db *DB) UpdateVineyard(ctx context.Context, vineyard *model.Vineyard) error {
	const query = `
    UPDATE vineyards 
    SET name = $1, location = $2 
    WHERE id = $3`
	_, err := db.ExecContext(ctx, query, vineyard.Name, vineyard.Location, vineyard.ID)
	if err != nil {
		return fmt.Errorf("updating vineyard: %w", err)
	}
	return nil
}

// DeleteVineyard removes a Vineyard record from the database.
func (db *DB) DeleteVineyard(ctx context.Context, id int) error {
	const query = `DELETE FROM vineyards WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting vineyard: %w", err)
	}
	return nil
}

// SaveSatelliteImageryMetadata stores metadata about satellite imagery for a vineyard.
func (db *DB) SaveSatelliteImageryMetadata(ctx context.Context, metadata *model.SatelliteData, vineyardID int) error {
	const query = `
    INSERT INTO satellite_imagery (vineyard_id, image_url, captured_at) 
    VALUES ($1, $2, $3)`
	_, err := db.ExecContext(ctx, query, vineyardID, metadata.ImageURL, metadata.CapturedAt)
	if err != nil {
		return fmt.Errorf("inserting satellite imagery metadata: %w", err)
	}
	return nil
}

// GetSatelliteImageryForVineyard retrieves all satellite imagery metadata for a specific vineyard.
func (db *DB) GetSatelliteImageryForVineyard(ctx context.Context, vineyardID int) ([]model.SatelliteData, error) {
	const query = `
    SELECT image_url, captured_at 
    FROM satellite_imagery 
    WHERE vineyard_id = $1`
	rows, err := db.QueryContext(ctx, query, vineyardID)
	if err != nil {
		return nil, fmt.Errorf("querying satellite imagery for vineyard: %w", err)
	}
	defer rows.Close()

	var images []model.SatelliteData
	for rows.Next() {
		var img model.SatelliteData
		if err := rows.Scan(&img.ImageURL, &img.CapturedAt); err != nil {
			return nil, fmt.Errorf("scanning satellite imagery: %w", err)
		}
		images = append(images, img)
	}
	return images, nil
}

// UpdateSoilData updates soil data for a given vineyard.
func (db *DB) UpdateSoilData(ctx context.Context, soilData *model.SoilData, vineyardID int) error {
	const query = `
    UPDATE vineyards 
    SET soil_health = $1 
    WHERE id = $2`
	// Assuming soil_health is a JSONB column or similar; can adapt as necessary
	soilHealthJSON, err := json.Marshal(soilData)
	if err != nil {
		return fmt.Errorf("marshaling soil data: %w", err)
	}

	_, err = db.ExecContext(ctx, query, soilHealthJSON, vineyardID)
	if err != nil {
		return fmt.Errorf("updating soil data: %w", err)
	}
	return nil
}

// GetVineyardWithEnvironmentalData retrieves a vineyard along with its related satellite imagery and soil data.
func (db *DB) GetVineyardWithEnvironmentalData(ctx context.Context, vineyardID int) (*model.Vineyard, error) {
	vineyard, err := db.GetVineyard(ctx, vineyardID)
	if err != nil {
		return nil, err
	}

	satelliteImagery, err := db.GetSatelliteImageryForVineyard(ctx, vineyardID)
	if err != nil {
		return nil, err
	}
	vineyard.SatelliteImagery = satelliteImagery

	// Assuming GetSoilDataForVineyard is not yet implemented
	soilData, err := db.GetSoilDataForVineyard(ctx, vineyardID)
	if err != nil {
		return nil, err
	}
	vineyard.SoilHealth = *soilData

	return vineyard, nil
}
