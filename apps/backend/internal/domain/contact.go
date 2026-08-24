package domain

import "time"

type ContactMessage struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(150)" json:"name"`
	Message      string    `gorm:"type:text" json:"message"`
	ContactEmail *string   `gorm:"type:varchar(255)" json:"contactEmail,omitempty"`
	ContactPhone *string   `gorm:"type:varchar(20)" json:"contactPhone,omitempty"`
	EmailSent    bool      `gorm:"default:false;index" json:"emailSent"`
	EmailError   string    `gorm:"type:text" json:"-"`
	CreatedAt    time.Time `gorm:"index" json:"createdAt"`
}

type ContactSubmission struct {
	Name         string
	Message      string
	ContactEmail string
	ContactPhone string
	Honeypot     string
}

type ValidationError struct{ Message string }

func (e *ValidationError) Error() string { return e.Message }

type MailMessage struct {
	Subject string
	Body    string
	ReplyTo string
}

type ContactRepository interface {
	Create(msg *ContactMessage) error
	UpdateEmailStatus(id uint, sent bool, errMsg string) error
}

type ContactUsecase interface {
	SubmitContactMessage(input ContactSubmission) error
}

type MailSender interface {
	Send(msg MailMessage) error
}
