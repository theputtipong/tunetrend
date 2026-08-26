package domain

// DeactivatedReasonAutoFetchFailure คือค่าที่ระบบตั้งเองเมื่อปิดหมวดหมู่อัตโนมัติ
const DeactivatedReasonAutoFetchFailure = "auto_fetch_failure"

type VideoCategory struct {
	ID          string `gorm:"primaryKey;type:varchar(10)" json:"id"`
	CountryCode string `gorm:"primaryKey;type:varchar(2)" json:"countryCode"`
	Title       string `gorm:"type:varchar(255)" json:"title"`
	Assignable  bool   `gorm:"index" json:"assignable"`
	// IsActive กำหนดว่าหมวดหมู่นี้แสดงให้ผู้ใช้เห็นหรือไม่
	IsActive          bool   `gorm:"not null;default:false;check:chk_video_categories_is_active_requires_assignable,is_active = false OR assignable = true" json:"isActive"`
	DeactivatedReason string `gorm:"type:varchar(50);default:''" json:"deactivatedReason,omitempty"`
	Note              string `gorm:"type:text" json:"note,omitempty"`
}

type VideoCategoryRepository interface {
	// UpsertCategoriesSetActive upsert title, assignable, is_active, deactivated_reason
	UpsertCategoriesSetActive(categories []VideoCategory) error
	// UpsertCategoriesPreserveActive upsert เฉพาะ title และ assignable
	UpsertCategoriesPreserveActive(categories []VideoCategory) error
	GetActiveCategories(countryCode string) ([]VideoCategory, error)
	// GetDeactivatedReasons คืนค่า deactivated_reason ปัจจุบันของแต่ละ id
	GetDeactivatedReasons(countryCode string, ids []string) (map[string]string, error)
}

type VideoCategoryUsecase interface {
	SyncCategories(countryCode string) error
	GetCategories(countryCode string) ([]VideoCategory, error)
}
