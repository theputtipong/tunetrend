package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"tunetrend-backend/internal/domain"
)

const resendAPIURL = "https://api.resend.com/emails"

const resendSandboxFrom = "onboarding@resend.dev"

var httpClient = &http.Client{Timeout: 10 * time.Second}

type ResendMailer struct {
	apiKey string
	to     string
}

func NewResendMailer() domain.MailSender {
	return &ResendMailer{
		apiKey: os.Getenv("RESEND_API_KEY"),
		to:     os.Getenv("CONTACT_EMAIL_TO"),
	}
}

type resendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
	ReplyTo string   `json:"reply_to,omitempty"`
}

type resendErrorResponse struct {
	Message string `json:"message"`
}

func (m *ResendMailer) Send(msg domain.MailMessage) error {
	payload := resendEmailRequest{
		From:    resendSandboxFrom,
		To:      []string{m.to},
		Subject: msg.Subject,
		Text:    msg.Body,
		ReplyTo: msg.ReplyTo,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("resend: ไม่สามารถ encode payload ได้: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, resendAPIURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("resend: สร้าง request ไม่สำเร็จ: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("resend: เรียก API ไม่สำเร็จ: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)

		var apiErr resendErrorResponse
		if jsonErr := json.Unmarshal(respBody, &apiErr); jsonErr == nil && apiErr.Message != "" {
			return fmt.Errorf("resend: api ตอบกลับ status %d: %s", resp.StatusCode, apiErr.Message)
		}
		return fmt.Errorf("resend: api ตอบกลับ status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
