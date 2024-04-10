type ImageService interface {
	ListImages(ctx context.Context, vineyardID string) ([]model.Image, error)
	SaveImage(ctx context.Context, image *model.Image) error
	GetImage(ctx context.Context, id string) (*model.Image, error)
	DeleteImage(ctx context.Context, id string) error
}

type imageServiceImpl struct {
	db      *db.DB
	storage *storage.StorageService
}

func NewImageService(db *db.DB, storage *storage.StorageService) ImageService {
	return &imageServiceImpl{db: db, storage: storage}
}

func (is *imageServiceImpl) ListImages(ctx context.Context, vineyardID string) ([]model.Image, error) {
	return is.db.GetSatelliteImageryForVineyard(ctx, vineyardID)
}

func (is *imageServiceImpl) SaveImage(ctx context.Context, image *model.Image) error {
	// Additional logic for storing the image file could be here if needed
	return is.db.SaveImage(ctx, image)
}

func (is *imageServiceImpl) GetImage(ctx context.Context, id string) (*model.Image, error) {
	return is.db.GetImage(ctx, id)
}

func (is *imageServiceImpl) DeleteImage(ctx context.Context, id string) error {
	return is.db.DeleteImage(ctx, id)
}
