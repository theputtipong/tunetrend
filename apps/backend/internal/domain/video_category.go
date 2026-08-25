package domain

// DeactivatedReasonAutoFetchFailure หมายถึง is_active=false ที่ระบบ sync ปิดให้เอง เพราะดึง content
// ไม่ได้ต่อเนื่องเกิน 24 ชม. — sync จะเปิดกลับให้อัตโนมัติทันทีที่ดึงสำเร็จอีกครั้ง
// ค่า DeactivatedReason อื่นใดที่ไม่ใช่ "" หรือค่านี้ ถือว่าเป็นการปิดด้วยเหตุผลอื่น (เช่นแอดมินตั้งเอง)
// sync จะไม่แตะ is_active ของหมวดนั้นอีกเลยจนกว่าจะมีคนเคลียร์ reason ออก
const DeactivatedReasonAutoFetchFailure = "auto_fetch_failure"

type VideoCategory struct {
	ID          string `gorm:"primaryKey;type:varchar(10)" json:"id"`
	CountryCode string `gorm:"primaryKey;type:varchar(2)" json:"countryCode"`
	Title       string `gorm:"type:varchar(255)" json:"title"`
	Assignable  bool   `gorm:"index" json:"assignable"`
	// IsActive คือสวิตช์ที่กำหนดว่าหมวดนี้พร้อมให้ผู้ใช้เลือกดูจริงหรือไม่ แยกจาก Assignable
	// (ค่าที่ YouTube เป็นคนบอกว่าหมวดนี้ "อนุญาตให้ผู้สร้างคลิปเลือกใช้" ในภูมิภาคนั้น) เพราะ assignable=true
	// ไม่ได้แปลว่าเราจะ sync/แสดงหมวดนั้นจริง (เช่นยังไม่พร้อม/ยังไม่ได้ทำตาราง sync/ดึง content ไม่ได้)
	// ตั้งแต่นี้ sync เป็นคนคำนวณค่านี้เองอัตโนมัติทุกครั้ง (ดู VideoCategoryUsecase.SyncCategories) ยกเว้น
	// แถวที่ DeactivatedReason ถูกตั้งเป็นค่าที่ไม่ใช่ auto (แปลว่ามีคนตั้งไว้เอง — sync จะไม่แตะอีก)
	// บังคับด้วย CHECK constraint ว่าถ้า assignable=false แล้ว is_active ต้องเป็น false เสมอ
	IsActive          bool   `gorm:"not null;default:false;check:chk_video_categories_is_active_requires_assignable,is_active = false OR assignable = true" json:"isActive"`
	DeactivatedReason string `gorm:"type:varchar(50);default:''" json:"deactivatedReason,omitempty"`
	Note              string `gorm:"type:text" json:"note,omitempty"`
}

type VideoCategoryRepository interface {
	// UpsertCategoriesSetActive upsert title/assignable/is_active/deactivated_reason —
	// ใช้ตอนที่ตัดสินใจค่า is_active ของแถวนั้นแล้วแน่นอน (ไม่ใช่ grace period/manual override)
	UpsertCategoriesSetActive(categories []VideoCategory) error
	// UpsertCategoriesPreserveActive upsert แค่ title/assignable — ไม่แตะ is_active/deactivated_reason เลย
	// ใช้ตอนยังอยู่ใน grace period (fail แต่ยังไม่ถึง 24 ชม.) หรือแถวถูก manual-override ไว้
	UpsertCategoriesPreserveActive(categories []VideoCategory) error
	GetActiveCategories(countryCode string) ([]VideoCategory, error)
	// GetDeactivatedReasons คืนค่า id -> deactivated_reason ปัจจุบันในตาราง สำหรับเช็คว่าหมวดไหน
	// ถูก manual-override ไว้อยู่ก่อนจะตัดสินใจว่าจะปรับ is_active ของ sync รอบนี้หรือไม่
	GetDeactivatedReasons(countryCode string, ids []string) (map[string]string, error)
}

type VideoCategoryUsecase interface {
	SyncCategories(countryCode string) error
	GetCategories(countryCode string) ([]VideoCategory, error)
}
