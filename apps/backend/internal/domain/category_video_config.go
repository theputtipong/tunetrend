package domain

type CategoryVideoConfig struct {
	CategoryID string `gorm:"primaryKey;type:varchar(10)" json:"categoryId"`
	TableName  string `gorm:"type:varchar(63);not null;unique" json:"tableName"`
	Label      string `gorm:"type:varchar(255);not null" json:"label"`
}

type CategoryVideoConfigRepository interface {
	GetAll() ([]CategoryVideoConfig, error)
}
