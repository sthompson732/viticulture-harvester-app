type SoilDataService interface {
	UpdateSoilData(ctx context.Context, vineyardID string, soilData *model.SoilData) error
	GetSoilData(ctx context.Context, vineyardID string) (*model.SoilData, error)
}

type soilDataServiceImpl struct {
	db *db.DB
}

func NewSoilDataService(db *db.DB) SoilDataService {
	return &soilDataServiceImpl{db: db}
}

func (sds *soilDataServiceImpl) UpdateSoilData(ctx context.Context, vineyardID string, soilData *model.SoilData) error {
	return sds.db.UpdateSoilData(ctx, vineyardID, soilData)
}

func (sds *soilDataServiceImpl) GetSoilData(ctx context.Context, vineyardID string) (*model.SoilData, error) {
	return sds.db.GetSoilDataForVineyard(ctx, vineyardID)
}