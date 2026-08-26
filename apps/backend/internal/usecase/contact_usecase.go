package usecase

import (
	"fmt"
	"log"

	"tunetrend-backend/internal/domain"
)

type contactUsecase struct {
	repo   domain.ContactRepository
	mailer domain.MailSender
}

func NewContactUsecase(repo domain.ContactRepository, mailer domain.MailSender) domain.ContactUsecase {
	return &contactUsecase{repo: repo, mailer: mailer}
}

func (u *contactUsecase) SubmitContactMessage(input domain.ContactSubmission) error {
	if input.Honeypot != "" {
		return nil
	}

	clean, err := validateAndSanitize(input)
	if err != nil {
		return err
	}

	msg := &domain.ContactMessage{
		Name:         clean.Name,
		Message:      clean.Message,
		ContactEmail: clean.ContactEmail,
		ContactPhone: clean.ContactPhone,
	}

	if err := u.repo.Create(msg); err != nil {
		return err
	}

	go func() {
		mailErr := u.mailer.Send(buildMailMessage(clean))
		if mailErr != nil {
			log.Printf("⚠️  [Contact] ส่งอีเมลแจ้งเตือนล้มเหลว (id=%d): %v", msg.ID, mailErr)
			_ = u.repo.UpdateEmailStatus(msg.ID, false, mailErr.Error())
			return
		}
		_ = u.repo.UpdateEmailStatus(msg.ID, true, "")
	}()

	return nil
}

func buildMailMessage(clean *cleanSubmission) domain.MailMessage {
	senderName := clean.Name
	if senderName == "" {
		senderName = "ไม่ระบุชื่อ"
	}

	contactLine := "ไม่ระบุช่องทางติดต่อกลับ"
	replyTo := ""
	if clean.ContactEmail != nil {
		contactLine = fmt.Sprintf("อีเมล: %s", *clean.ContactEmail)
		replyTo = *clean.ContactEmail
	} else if clean.ContactPhone != nil {
		contactLine = fmt.Sprintf("เบอร์โทร: %s", *clean.ContactPhone)
	}

	body := fmt.Sprintf(
		"มีข้อความใหม่จากฟอร์ม Contact ของ TuneTrend\n\nชื่อ: %s\n%s\n\nข้อความ:\n%s\n",
		senderName, contactLine, clean.Message,
	)

	return domain.MailMessage{
		Subject: fmt.Sprintf("[TuneTrend Contact] ข้อความจาก %s", sanitizeHeaderValue(senderName)),
		Body:    body,
		ReplyTo: replyTo,
	}
}
