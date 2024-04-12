package service

import (
	"context"
	"errors"
	"time"

	"github.com/sthompson732/viticulture-harvester-app/internal/db"
	"github.com/sthompson732/viticulture-harvester-app/internal/model"
)

type WeatherService interface {
	CreateWeatherData(ctx context.Context, weather *model.WeatherData) error
	GetWeatherData(ctx context.Context, id int) (*model.WeatherData, error)
	UpdateWeatherData(ctx context.Context, weather *model.WeatherData) error
	DeleteWeatherData(ctx context.Context, id int) error
	ListWeatherDataByVineyard(ctx context.Context, vineyardID int) ([]model.WeatherData, error)
	ListWeatherDataByDateRange(ctx context.Context, vineyardID int, start, end time.Time) ([]model.WeatherData, error)
}

type weatherServiceImpl struct {
	db *db.DB
}

func NewWeatherService(db *db.DB) WeatherService {
	return &weatherServiceImpl{db: db}
}

func (ws *weatherServiceImpl) CreateWeatherData(ctx context.Context, weather *model.WeatherData) error {
	if weather == nil {
		return errors.New("cannot create nil weather data")
	}
	return ws.db.SaveWeatherData(ctx, weather)
}

func (ws *weatherServiceImpl) GetWeatherData(ctx context.Context, id int) (*model.WeatherData, error) {
	if id <= 0 {
		return nil, errors.New("invalid weather data ID")
	}
	return ws.db.GetWeatherData(ctx, id)
}

func (ws *weatherServiceImpl) UpdateWeatherData(ctx context.Context, weather *model.WeatherData) error {
	if weather == nil {
		return errors.New("cannot update nil weather data")
	}
	if weather.ID == 0 {
		return errors.New("invalid weather data ID")
	}
	return ws.db.UpdateWeatherData(ctx, weather)
}

func (ws *weatherServiceImpl) DeleteWeatherData(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid weather data ID")
	}
	return ws.db.DeleteWeatherData(ctx, id)
}

func (ws *weatherServiceImpl) ListWeatherDataByVineyard(ctx context.Context, vineyardID int) ([]model.WeatherData, error) {
	if vineyardID <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	return ws.db.ListWeatherDataByVineyard(ctx, vineyardID)
}

func (ws *weatherServiceImpl) ListWeatherDataByDateRange(ctx context.Context, vineyardID int, start, end time.Time) ([]model.WeatherData, error) {
	if vineyardID <= 0 {
		return nil, errors.New("invalid vineyard ID")
	}
	if start.After(end) {
		return nil, errors.New("start date must be before end date")
	}
	return ws.db.ListWeatherDataByDateRange(ctx, vineyardID, start, end)
}
