package database

import (
	"tunetrend-backend/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var defaultCategoryVideoConfigs = []domain.CategoryVideoConfig{
	{CategoryID: "1", TableName: "film_animation_videos", Label: "Film & Animation"},
	{CategoryID: "2", TableName: "autos_vehicles_videos", Label: "Autos & Vehicles"},
	{CategoryID: "15", TableName: "pets_animals_videos", Label: "Pets & Animals"},
	{CategoryID: "17", TableName: "sports_videos", Label: "Sports"},
	{CategoryID: "19", TableName: "travel_events_videos", Label: "Travel & Events"},
	{CategoryID: "20", TableName: "gaming_videos", Label: "Gaming"},
	{CategoryID: "22", TableName: "people_blogs_videos", Label: "People & Blogs"},
	{CategoryID: "23", TableName: "comedy_videos", Label: "Comedy"},
	{CategoryID: "24", TableName: "entertainment_videos", Label: "Entertainment"},
	{CategoryID: "25", TableName: "news_politics_videos", Label: "News & Politics"},
	{CategoryID: "26", TableName: "howto_style_videos", Label: "Howto & Style"},
	{CategoryID: "27", TableName: "education_videos", Label: "Education"},
	{CategoryID: "28", TableName: "science_technology_videos", Label: "Science & Technology"},
}

const defaultCountries = "TH,KR,JP,US,GB"

var defaultWorkerSettings = domain.WorkerSettings{
	ID:                               1,
	Countries:                        defaultCountries,
	MusicSyncIntervalMinutes:         180,
	CategoryVideoSyncIntervalMinutes: 180,
	CategoryResumeIntervalMinutes:    10080,
	CategorySyncIntervalMinutes:      1440,
	CategoryFetchFailureGraceMinutes: 1440,
}

func seedCategoryVideoConfigs(db *gorm.DB) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "category_id"}},
		DoNothing: true,
	}).Create(&defaultCategoryVideoConfigs).Error
}

func seedWorkerSettings(db *gorm.DB) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&defaultWorkerSettings).Error
}
