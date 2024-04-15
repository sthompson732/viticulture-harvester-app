/*
 * db.go: Database interaction layer for the app.
 * Contains methods for executing CRUD operations on the database.
 * Usage: Utilized by service layers to access and manipulate database records.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */

package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/sthompson732/viticulture-harvester-app/internal/model"
)

type DB struct {
	*sql.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error verifying connection to database: %w", err)
	}

	return &DB{db}, nil
}

// Image methods
func (db *DB) SaveImage(ctx context.Context, image *model.Image) error {
	const query = `
    INSERT INTO images (vineyard_id, url, captured_at)
    VALUES ($1, $2, $3)
    RETURNING id`
	err := db.QueryRowContext(ctx, query, image.VineyardID, image.URL, image.CapturedAt).Scan(&image.ID)
	if err != nil {
		return fmt.Errorf("inserting image: %w", err)
	}
	return nil
}

func (db *DB) GetImage(ctx context.Context, id int) (*model.Image, error) {
	const query = `
    SELECT id, vineyard_id, url, captured_at
    FROM images
    WHERE id = $1`
	img := &model.Image{}
	err := db.QueryRowContext(ctx, query, id).Scan(&img.ID, &img.VineyardID, &img.URL, &img.CapturedAt)
	if err != nil {
		return nil, fmt.Errorf("retrieving image by ID: %w", err)
	}
	return img, nil
}

func (db *DB) DeleteImage(ctx context.Context, id int) error {
	const query = `DELETE FROM images WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}
	return nil
}

func (db *DB) FindImagesByDateRange(ctx context.Context, vineyardID int, start, end time.Time) ([]model.Image, error) {
	query := `SELECT id, vineyard_id, image_url, description, captured_at, bbox FROM images
              WHERE vineyard_id = $1 AND captured_at BETWEEN $2 AND $3`
	rows, err := db.QueryContext(ctx, query, vineyardID, start, end)
	if err != nil {
		return nil, fmt.Errorf("error querying images by date range: %w", err)
	}
	defer rows.Close()

	var images []model.Image
	for rows.Next() {
		var img model.Image
		if err := rows.Scan(&img.ID, &img.VineyardID, &img.URL, &img.Description, &img.CapturedAt, &img.BoundingBox); err != nil {
			return nil, fmt.Errorf("error scanning image: %w", err)
		}
		images = append(images, img)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading images: %w", err)
	}
	return images, nil
}

func (db *DB) GetRecentImages(ctx context.Context, vineyardID int, limit int) ([]model.Image, error) {
	query := `SELECT id, vineyard_id, image_url, description, captured_at, bbox FROM images
              WHERE vineyard_id = $1 ORDER BY captured_at DESC LIMIT $2`
	rows, err := db.QueryContext(ctx, query, vineyardID, limit)
	if err != nil {
		return nil, fmt.Errorf("error querying recent images: %w", err)
	}
	defer rows.Close()

	var images []model.Image
	for rows.Next() {
		var img model.Image
		if err := rows.Scan(&img.ID, &img.VineyardID, &img.URL, &img.Description, &img.CapturedAt, &img.BoundingBox); err != nil {
			return nil, fmt.Errorf("error scanning image: %w", err)
		}
		images = append(images, img)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading images: %w", err)
	}
	return images, nil
}

// ListImagesByVineyard retrieves all images for a specific vineyard.
func (db *DB) ListImagesByVineyard(ctx context.Context, vineyardID int) ([]model.Image, error) {
	const query = `
    SELECT id, vineyard_id, image_url, description, captured_at, bbox
    FROM images
    WHERE vineyard_id = $1`
	rows, err := db.QueryContext(ctx, query, vineyardID)
	if err != nil {
		return nil, fmt.Errorf("querying images for vineyard: %w", err)
	}
	defer rows.Close()

	var images []model.Image
	for rows.Next() {
		var img model.Image
		err := rows.Scan(&img.ID, &img.VineyardID, &img.URL, &img.Description, &img.CapturedAt, &img.BoundingBox)
		if err != nil {
			return nil, fmt.Errorf("scanning image: %w", err)
		}
		images = append(images, img)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("reading image rows: %w", err)
	}
	return images, nil
}

// UpdateImage updates the details for an existing image.
func (db *DB) UpdateImage(ctx context.Context, image *model.Image) error {
	const query = `
    UPDATE images
    SET vineyard_id = $1, image_url = $2, description = $3, captured_at = $4, bbox = $5
    WHERE id = $6`
	_, err := db.ExecContext(ctx, query, image.VineyardID, image.URL, image.Description, image.CapturedAt, image.BoundingBox, image.ID)
	if err != nil {
		return fmt.Errorf("updating image: %w", err)
	}
	return nil
}

// Vineyard methods
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

// ListVineyards retrieves all vineyard entries from the database.
func (db *DB) ListVineyards(ctx context.Context) ([]model.Vineyard, error) {
	const query = `
    SELECT id, name, location, ST_AsText(bbox) AS bbox_text
    FROM vineyards`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying vineyards: %w", err)
	}
	defer rows.Close()

	var vineyards []model.Vineyard
	for rows.Next() {
		var vineyard model.Vineyard
		var bboxText string // We use this to hold the bbox polygon text
		err := rows.Scan(&vineyard.ID, &vineyard.Name, &vineyard.Location, &bboxText)
		if err != nil {
			return nil, fmt.Errorf("scanning vineyard: %w", err)
		}
		// Convert bbox text back to polygon type if necessary, this step may require further parsing depending on how you handle geometries
		vineyard.BoundingBox = bboxText
		vineyards = append(vineyards, vineyard)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("reading vineyard rows: %w", err)
	}
	return vineyards, nil
}

// Satellite Imagery methods
// SaveSatelliteImagery stores new satellite imagery data.
func (db *DB) SaveSatelliteImagery(ctx context.Context, sd *model.SatelliteData) error {
	query := `
    INSERT INTO satellite_imagery (vineyard_id, image_url, captured_at, bbox)
    VALUES ($1, $2, $3, $4)
    RETURNING id`
	err := db.QueryRowContext(ctx, query, sd.VineyardID, sd.ImageURL, sd.CapturedAt, sd.BoundingBox).Scan(&sd.ID)
	if err != nil {
		return fmt.Errorf("error inserting satellite imagery: %w", err)
	}
	return nil
}

// GetSatelliteImagery retrieves a single satellite imagery record by ID.
func (db *DB) GetSatelliteImagery(ctx context.Context, id int) (*model.SatelliteData, error) {
	query := `
    SELECT id, vineyard_id, image_url, captured_at, bbox
    FROM satellite_imagery
    WHERE id = $1`
	var sd model.SatelliteData
	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(&sd.ID, &sd.VineyardID, &sd.ImageURL, &sd.CapturedAt, &sd.BoundingBox)
	if err != nil {
		return nil, fmt.Errorf("error retrieving satellite imagery: %w", err)
	}
	return &sd, nil
}

// DeleteSatelliteImagery deletes a satellite imagery record by ID.
func (db *DB) DeleteSatelliteImagery(ctx context.Context, id int) error {
	query := `DELETE FROM satellite_imagery WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting satellite imagery: %w", err)
	}
	return nil
}

// UpdateSatelliteImagery updates an existing satellite imagery record.
func (db *DB) UpdateSatelliteImagery(ctx context.Context, sd *model.SatelliteData) error {
	query := `
    UPDATE satellite_imagery
    SET image_url = $1, captured_at = $2, bbox = $3, vineyard_id = $4
    WHERE id = $5`
	_, err := db.ExecContext(ctx, query, sd.ImageURL, sd.CapturedAt, sd.BoundingBox, sd.VineyardID, sd.ID)
	if err != nil {
		return fmt.Errorf("error updating satellite imagery: %w", err)
	}
	return nil
}

// SaveSatelliteImageryMetadata stores metadata about satellite imagery for a vineyard.
func (db *DB) SaveSatelliteImageryMetadata(ctx context.Context, data *model.SatelliteData, vineyardID int) error {
	// SQL execution logic here, for example:
	const query = `INSERT INTO satellite_imagery (vineyard_id, image_url, resolution, captured_at, bounding_box)
                   VALUES ($1, $2, $3, $4, $5)`
	_, err := db.ExecContext(ctx, query, vineyardID, data.ImageURL, data.Resolution, data.CapturedAt, data.BoundingBox)
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

// ListSatelliteImageryByVineyard retrieves all satellite imagery for a specific vineyard.
func (db *DB) ListSatelliteImageryByVineyard(ctx context.Context, vineyardID int) ([]model.SatelliteData, error) {
	const query = `
    SELECT id, vineyard_id, image_url, resolution, captured_at, bbox
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
		err := rows.Scan(&img.ID, &img.VineyardID, &img.ImageURL, &img.Resolution, &img.CapturedAt, &img.BoundingBox)
		if err != nil {
			return nil, fmt.Errorf("scanning satellite imagery: %w", err)
		}
		images = append(images, img)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("reading satellite imagery rows: %w", err)
	}
	return images, nil
}

// ListSatelliteImageryByDateRange retrieves satellite imagery within a specified date range for a vineyard.
func (db *DB) ListSatelliteImageryByDateRange(ctx context.Context, vineyardID int, startDate, endDate time.Time) ([]model.SatelliteData, error) {
	const query = `
    SELECT id, vineyard_id, image_url, resolution, captured_at, ST_AsText(bbox) AS bbox_text
    FROM satellite_imagery
    WHERE vineyard_id = $1 AND captured_at BETWEEN $2 AND $3`

	rows, err := db.QueryContext(ctx, query, vineyardID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("querying satellite imagery by date range: %w", err)
	}
	defer rows.Close()

	var images []model.SatelliteData
	for rows.Next() {
		var img model.SatelliteData
		var bboxText string
		err := rows.Scan(&img.ID, &img.VineyardID, &img.ImageURL, &img.Resolution, &img.CapturedAt, &bboxText)
		if err != nil {
			return nil, fmt.Errorf("scanning satellite imagery: %w", err)
		}
		img.BoundingBox = bboxText // Assuming bbox is needed as text; convert if necessary to your preferred format
		images = append(images, img)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("reading satellite imagery rows: %w", err)
	}
	return images, nil
}

// Soil methods
// SaveSoilData inserts a new SoilData record into the database.
func (db *DB) SaveSoilData(ctx context.Context, soilData *model.SoilData) error {
	const query = `
    INSERT INTO soil_data (vineyard_id, data, location, sampled_at)
    VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326), $5)
    RETURNING id`
	jsonData, err := json.Marshal(soilData)
	if err != nil {
		return fmt.Errorf("error marshaling soil data: %w", err)
	}
	err = db.QueryRowContext(ctx, query, soilData.VineyardID, jsonData, soilData.Location.X, soilData.Location.Y, soilData.SampledAt).Scan(&soilData.ID)
	if err != nil {
		return fmt.Errorf("error inserting soil data: %w", err)
	}
	return nil
}

// DeleteSoilData removes a SoilData record from the database by ID.
func (db *DB) DeleteSoilData(ctx context.Context, id int) error {
	const query = `DELETE FROM soil_data WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting soil data: %w", err)
	}
	return nil
}

// UpdateSoilData updates a given SoilData's details.
func (db *DB) UpdateSoilData(ctx context.Context, soilData *model.SoilData) error {
	jsonData, err := json.Marshal(soilData)
	if err != nil {
		return fmt.Errorf("error marshaling soil data: %w", err)
	}

	query := `
    UPDATE soil_data
    SET data = $1, location = ST_SetSRID(ST_MakePoint($2, $3), 4326), sampled_at = $4
    WHERE id = $5`

	_, err = db.ExecContext(ctx, query, jsonData, soilData.Location.X, soilData.Location.Y, soilData.SampledAt, soilData.ID)
	if err != nil {
		return fmt.Errorf("updating soil data: %w", err)
	}
	return nil
}

// ListSoilDataForVineyard retrieves all SoilData for a specific vineyard.
func (db *DB) ListSoilDataForVineyard(ctx context.Context, vineyardID int) ([]model.SoilData, error) {
	const query = `
    SELECT id, vineyard_id, data, ST_X(location) AS longitude, ST_Y(location) AS latitude, sampled_at
    FROM soil_data
    WHERE vineyard_id = $1`

	rows, err := db.QueryContext(ctx, query, vineyardID)
	if err != nil {
		return nil, fmt.Errorf("querying soil data for vineyard: %w", err)
	}
	defer rows.Close()

	var soils []model.SoilData
	for rows.Next() {
		var soil model.SoilData
		var jsonData []byte
		if err := rows.Scan(&soil.ID, &soil.VineyardID, &jsonData, &soil.Location.X, &soil.Location.Y, &soil.SampledAt); err != nil {
			return nil, fmt.Errorf("scanning soil data: %w", err)
		}
		if err = json.Unmarshal(jsonData, &soil); err != nil {
			return nil, fmt.Errorf("unmarshaling soil data: %w", err)
		}
		soils = append(soils, soil)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("reading soil data rows: %w", err)
	}
	return soils, nil
}

// ListSoilDataByDateRange retrieves soil data within a specified date range for a vineyard.
func (db *DB) ListSoilDataByDateRange(ctx context.Context, vineyardID int, start, end time.Time) ([]model.SoilData, error) {
	const query = `
    SELECT id, vineyard_id, data, ST_X(location) AS longitude, ST_Y(location) AS latitude, sampled_at
    FROM soil_data
    WHERE vineyard_id = $1 AND sampled_at BETWEEN $2 AND $3`

	rows, err := db.QueryContext(ctx, query, vineyardID, start, end)
	if err != nil {
		return nil, fmt.Errorf("querying soil data by date range: %w", err)
	}
	defer rows.Close()

	var soils []model.SoilData
	for rows.Next() {
		var soil model.SoilData
		var jsonData []byte
		if err := rows.Scan(&soil.ID, &soil.VineyardID, &jsonData, &soil.Location.X, &soil.Location.Y, &soil.SampledAt); err != nil {
			return nil, fmt.Errorf("scanning soil data: %w", err)
		}
		if err = json.Unmarshal(jsonData, &soil); err != nil {
			return nil, fmt.Errorf("unmarshaling soil data: %w", err)
		}
		soils = append(soils, soil)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("reading soil data rows: %w", err)
	}
	return soils, nil
}

// GetSoilData retrieves SoilData by ID.
func (db *DB) GetSoilData(ctx context.Context, id int) (*model.SoilData, error) {
	const query = `
    SELECT id, vineyard_id, data, ST_X(location) AS longitude, ST_Y(location) AS latitude, sampled_at
    FROM soil_data
    WHERE id = $1`
	soilData := &model.SoilData{}
	row := db.QueryRowContext(ctx, query, id)
	var jsonData []byte
	err := row.Scan(&soilData.ID, &soilData.VineyardID, &jsonData, &soilData.Location.X, &soilData.Location.Y, &soilData.SampledAt)
	if err != nil {
		return nil, fmt.Errorf("retrieving soil data by ID: %w", err)
	}
	if err = json.Unmarshal(jsonData, &soilData); err != nil {
		return nil, fmt.Errorf("unmarshaling soil data: %w", err)
	}
	return soilData, nil
}

// GetVineyardWithEnvironmentalData retrieves a vineyard along with its related satellite imagery and soil data.
func (db *DB) GetVineyardWithEnvironmentalData(ctx context.Context, vineyardID int) (*model.Vineyard, error) {
	vineyard, err := db.GetVineyard(ctx, vineyardID)
	if err != nil {
		return nil, fmt.Errorf("retrieving vineyard by ID: %w", err)
	}

	satelliteImagery, err := db.GetSatelliteImageryForVineyard(ctx, vineyardID)
	if err != nil {
		return nil, fmt.Errorf("retrieving satellite imagery for vineyard: %w", err)
	}
	vineyard.SatelliteImagery = satelliteImagery

	soilData, err := db.GetSoilDataForVineyard(ctx, vineyardID)
	if err != nil {
		return nil, fmt.Errorf("retrieving soil data for vineyard: %w", err)
	}
	vineyard.SoilHealth = soilData

	return vineyard, nil
}

// GetSoilDataForVineyard retrieves all soil data entries for a specific vineyard.
func (db *DB) GetSoilDataForVineyard(ctx context.Context, vineyardID int) ([]model.SoilData, error) {
	const query = `
    SELECT data
    FROM soil_data
    WHERE vineyard_id = $1`

	rows, err := db.QueryContext(ctx, query, vineyardID)
	if err != nil {
		return nil, fmt.Errorf("querying soil data for vineyard: %w", err)
	}
	defer rows.Close()

	var soils []model.SoilData
	for rows.Next() {
		var data model.SoilData
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("scanning soil data: %w", err)
		}
		soils = append(soils, data)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("reading soil data rows: %w", err)
	}

	return soils, nil
}

// Pest methods
// SavePestData inserts a new PestData record into the database.
func (db *DB) SavePestData(ctx context.Context, pest *model.PestData) error {
	const query = `
    INSERT INTO pest_data (vineyard_id, description, observation_date, location, pest_type, severity)
    VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5), 4326), $6, $7)
    RETURNING id`
	err := db.QueryRowContext(ctx, query, pest.VineyardID, pest.Description, pest.ObservationDate, pest.Location.X, pest.Location.Y, pest.PestType, pest.Severity).Scan(&pest.ID)
	if err != nil {
		return fmt.Errorf("inserting pest data: %w", err)
	}
	return nil
}

// GetPestData retrieves a PestData by ID.
func (db *DB) GetPestData(ctx context.Context, id int) (*model.PestData, error) {
	const query = `
    SELECT id, vineyard_id, description, observation_date, ST_X(location) AS longitude, ST_Y(location) AS latitude, pest_type, severity
    FROM pest_data
    WHERE id = $1`
	pest := &model.PestData{}
	err := db.QueryRowContext(ctx, query, id).Scan(&pest.ID, &pest.VineyardID, &pest.Description, &pest.ObservationDate, &pest.Location.X, &pest.Location.Y, &pest.Type, &pest.Severity)
	if err != nil {
		return nil, fmt.Errorf("retrieving pest data by ID: %w", err)
	}
	return pest, nil
}

// UpdatePestData updates a given PestData's details.
func (db *DB) UpdatePestData(ctx context.Context, pest *model.PestData) error {
	const query = `
    UPDATE pest_data
    SET description = $1, observation_date = $2, location = ST_SetSRID(ST_MakePoint($3, $4), 4326), pest_type = $5, severity = $6
    WHERE id = $7`
	_, err := db.ExecContext(ctx, query, pest.Description, pest.ObservationDate, pest.Location.X, pest.Location.Y, pest.Type, pest.Severity, pest.ID)
	if err != nil {
		return fmt.Errorf("updating pest data: %w", err)
	}
	return nil
}

// DeletePestData removes a PestData record from the database.
func (db *DB) DeletePestData(ctx context.Context, id int) error {
	const query = `DELETE FROM pest_data WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting pest data: %w", err)
	}
	return nil
}

// ListPestDataByVineyard retrieves all PestData for a specific vineyard.
func (db *DB) ListPestDataByVineyard(ctx context.Context, vineyardID int) ([]model.PestData, error) {
	const query = `
    SELECT id, vineyard_id, description, observation_date, ST_X(location) AS longitude, ST_Y(location) AS latitude, pest_type, severity
    FROM pest_data
    WHERE vineyard_id = $1`
	rows, err := db.QueryContext(ctx, query, vineyardID)
	if err != nil {
		return nil, fmt.Errorf("querying pest data for vineyard: %w", err)
	}
	defer rows.Close()

	var pests []model.PestData
	for rows.Next() {
		var pest model.PestData
		if err := rows.Scan(&pest.ID, &pest.VineyardID, &pest.Description, &pest.ObservationDate, &pest.Location.X, &pest.Location.Y, &pest.Type, &pest.Severity); err != nil {
			return nil, fmt.Errorf("scanning pest data: %w", err)
		}
		pests = append(pests, pest)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("reading pest data rows: %w", err)
	}
	return pests, nil
}

// ListPestDataByDateRange retrieves PestData for a specific vineyard within a date range.
func (db *DB) ListPestDataByDateRange(ctx context.Context, vineyardID int, start, end time.Time) ([]model.PestData, error) {
	const query = `
    SELECT id, vineyard_id, description, observation_date, ST_X(location) AS longitude, ST_Y(location) AS latitude, pest_type, severity
    FROM pest_data
    WHERE vineyard_id = $1 AND observation_date BETWEEN $2 AND $3`
	rows, err := db.QueryContext(ctx, query, vineyardID, start, end)
	if err != nil {
		return nil, fmt.Errorf("querying pest data by date range: %w", err)
	}
	defer rows.Close()

	var pests []model.PestData
	for rows.Next() {
		var pest model.PestData
		if err := rows.Scan(&pest.ID, &pest.VineyardID, &pest.Description, &pest.ObservationDate, &pest.Location.X, &pest.Location.Y, &pest.Type, &pest.Severity); err != nil {
			return nil, fmt.Errorf("scanning pest data: %w", err)
		}
		pests = append(pests, pest)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("reading pest data rows: %w", err)
	}
	return pests, nil
}

func (db *DB) FilterPestData(ctx context.Context, vineyardID int, pestType, severity string) ([]model.PestData, error) {
	query := `SELECT id, vineyard_id, description, observation_date, location, pest_type, severity FROM pest_data
              WHERE vineyard_id = $1 AND pest_type = $2 AND severity = $3`
	rows, err := db.QueryContext(ctx, query, vineyardID, pestType, severity)
	if err != nil {
		return nil, fmt.Errorf("error querying filtered pest data: %w", err)
	}
	defer rows.Close()

	var pests []model.PestData
	for rows.Next() {
		var pest model.PestData
		err := rows.Scan(&pest.ID, &pest.VineyardID, &pest.Description, &pest.ObservationDate, &pest.Location, &pest.Type, &pest.Severity)
		if err != nil {
			return nil, fmt.Errorf("error scanning pest data: %w", err)
		}
		pests = append(pests, pest)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading pest data rows: %w", err)
	}
	return pests, nil
}

// Weather methods
// SaveWeatherData inserts a new WeatherData record into the database.
func (db *DB) SaveWeatherData(ctx context.Context, weather *model.WeatherData) error {
	const query = `
    INSERT INTO weather_data (vineyard_id, temperature, humidity, observation_time, location)
    VALUES ($1, $2, $3, $4, ST_SetSRID(ST_MakePoint($5, $6), 4326))
    RETURNING id`
	err := db.QueryRowContext(ctx, query, weather.VineyardID, weather.Temperature, weather.Humidity, weather.ObservationTime, weather.Location.X, weather.Location.Y).Scan(&weather.ID)
	if err != nil {
		return fmt.Errorf("inserting weather data: %w", err)
	}
	return nil
}

// GetWeatherData retrieves a WeatherData by ID.
func (db *DB) GetWeatherData(ctx context.Context, id int) (*model.WeatherData, error) {
	const query = `
    SELECT id, vineyard_id, temperature, humidity, observation_time, ST_X(location) AS longitude, ST_Y(location) AS latitude
    FROM weather_data
    WHERE id = $1`
	weather := &model.WeatherData{}
	err := db.QueryRowContext(ctx, query, id).Scan(&weather.ID, &weather.VineyardID, &weather.Temperature, &weather.Humidity, &weather.ObservationTime, &weather.Location.X, &weather.Location.Y)
	if err != nil {
		return nil, fmt.Errorf("retrieving weather data by ID: %w", err)
	}
	return weather, nil
}

// UpdateWeatherData updates a given WeatherData's details.
func (db *DB) UpdateWeatherData(ctx context.Context, weather *model.WeatherData) error {
	const query = `
    UPDATE weather_data
    SET temperature = $1, humidity = $2, observation_time = $3, location = ST_SetSRID(ST_MakePoint($4, $5), 4326)
    WHERE id = $6`
	_, err := db.ExecContext(ctx, query, weather.Temperature, weather.Humidity, weather.ObservationTime, weather.Location.X, weather.Location.Y, weather.ID)
	if err != nil {
		return fmt.Errorf("updating weather data: %w", err)
	}
	return nil
}

// DeleteWeatherData removes a WeatherData record from the database.
func (db *DB) DeleteWeatherData(ctx context.Context, id int) error {
	const query = `DELETE FROM weather_data WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting weather data: %w", err)
	}
	return nil
}

// ListWeatherDataByVineyard retrieves all WeatherData for a specific vineyard.
func (db *DB) ListWeatherDataByVineyard(ctx context.Context, vineyardID int) ([]model.WeatherData, error) {
	const query = `
    SELECT id, vineyard_id, temperature, humidity, observation_time, ST_X(location) AS longitude, ST_Y(location) AS latitude
    FROM weather_data
    WHERE vineyard_id = $1`
	rows, err := db.QueryContext(ctx, query, vineyardID)
	if err != nil {
		return nil, fmt.Errorf("querying weather data for vineyard: %w", err)
	}
	defer rows.Close()

	var weathers []model.WeatherData
	for rows.Next() {
		var weather model.WeatherData
		if err := rows.Scan(&weather.ID, &weather.VineyardID, &weather.Temperature, &weather.Humidity, &weather.ObservationTime, &weather.Location.X, &weather.Location.Y); err != nil {
			return nil, fmt.Errorf("scanning weather data: %w", err)
		}
		weathers = append(weathers, weather)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("reading weather data rows: %w", err)
	}
	return weathers, nil
}

// ListWeatherDataByDateRange retrieves WeatherData for a specific vineyard within a date range.
func (db *DB) ListWeatherDataByDateRange(ctx context.Context, vineyardID int, start, end time.Time) ([]model.WeatherData, error) {
	const query = `
    SELECT id, vineyard_id, temperature, humidity, observation_time, ST_X(location) AS longitude, ST_Y(location) AS latitude
    FROM weather_data
    WHERE vineyard_id = $1 AND observation_time BETWEEN $2 AND $3`
	rows, err := db.QueryContext(ctx, query, vineyardID, start, end)
	if err != nil {
		return nil, fmt.Errorf("querying weather data by date range: %w", err)
	}
	defer rows.Close()

	var weathers []model.WeatherData
	for rows.Next() {
		var weather model.WeatherData
		if err := rows.Scan(&weather.ID, &weather.VineyardID, &weather.Temperature, &weather.Humidity, &weather.ObservationTime, &weather.Location.X, &weather.Location.Y); err != nil {
			return nil, fmt.Errorf("scanning weather data: %w", err)
		}
		weathers = append(weathers, weather)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("reading weather data rows: %w", err)
	}
	return weathers, nil
}
