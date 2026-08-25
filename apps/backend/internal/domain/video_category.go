package domain

type VideoCategory struct {
	ID          string `gorm:"primaryKey;type:varchar(10)" json:"id"`
	CountryCode string `gorm:"primaryKey;type:varchar(2)" json:"countryCode"`
	Title       string `gorm:"type:varchar(255)" json:"title"`
	Assignable  bool   `gorm:"index" json:"assignable"`
	// IsActive คือสวิตช์ที่แอดมินเป็นคนกำหนดเองว่าหมวดนี้พร้อมให้ผู้ใช้เลือกดูจริงหรือไม่
	// แยกจาก Assignable (ค่าที่ YouTube เป็นคนบอกว่าหมวดนี้ "อนุญาตให้ผู้สร้างคลิปเลือกใช้" ในภูมิภาคนั้น)
	// เพราะ assignable=true ไม่ได้แปลว่าเราจะ sync/แสดงหมวดนั้นจริง (เช่นยังไม่พร้อม/ยังไม่ได้ทำตาราง sync)
	// บังคับด้วย CHECK constraint ว่าถ้า assignable=false แล้ว is_active ต้องเป็น false เสมอ
	IsActive bool   `gorm:"not null;default:false;check:chk_video_categories_is_active_requires_assignable,is_active = false OR assignable = true" json:"isActive"`
	Note     string `gorm:"type:text" json:"note,omitempty"`
}

type VideoCategoryRepository interface {
	UpsertCategories(categories []VideoCategory) error
	GetActiveCategories(countryCode string) ([]VideoCategory, error)
}

type VideoCategoryUsecase interface {
	SyncCategories(countryCode string) error
	GetCategories(countryCode string) ([]VideoCategory, error)
}
