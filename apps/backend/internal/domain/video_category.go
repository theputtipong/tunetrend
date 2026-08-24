package domain

type VideoCategory struct {
	ID          string `gorm:"primaryKey;type:varchar(10)" json:"id"`
	CountryCode string `gorm:"primaryKey;type:varchar(2)" json:"countryCode"`
	Title       string `gorm:"type:varchar(255)" json:"title"`
	Assignable  bool   `gorm:"index" json:"assignable"`
}

type VideoCategoryRepository interface {
	UpsertCategories(categories []VideoCategory) error
	GetAssignableCategories(countryCode string) ([]VideoCategory, error)
}

type VideoCategoryUsecase interface {
	SyncCategories(countryCode string) error
	GetCategories(countryCode string) ([]VideoCategory, error)
}
