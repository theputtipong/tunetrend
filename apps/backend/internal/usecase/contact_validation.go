package usecase

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"

	"tunetrend-backend/internal/domain"
)

const (
	minMessageRunes = 10
	maxMessageBytes = 5000
	maxNameBytes    = 150
)

var thaiPhoneRegex = regexp.MustCompile(`^0(?:[689]\d{8}|[2-57]\d{7})$`)

func normalizePhone(raw string) string {
	replacer := strings.NewReplacer(" ", "", "-", "", "(", "", ")", "")
	return replacer.Replace(strings.TrimSpace(raw))
}

func isValidThaiPhone(raw string) bool {
	return thaiPhoneRegex.MatchString(normalizePhone(raw))
}

func isValidEmail(raw string) bool {
	trimmed := strings.TrimSpace(raw)
	addr, err := mail.ParseAddress(trimmed)
	if err != nil {
		return false
	}
	return addr.Address == trimmed
}

func sanitizeHeaderValue(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	return strings.TrimSpace(s)
}

func stripControlChars(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\t' || (r >= 0x20 && r != 0x7f) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

type cleanSubmission struct {
	Name         string
	Message      string
	ContactEmail *string
	ContactPhone *string
}

func validateAndSanitize(input domain.ContactSubmission) (*cleanSubmission, error) {
	name := sanitizeHeaderValue(stripControlChars(input.Name))
	if len(name) > maxNameBytes {
		return nil, &domain.ValidationError{Message: "name is too long"}
	}

	message := strings.TrimSpace(stripControlChars(input.Message))
	if utf8.RuneCountInString(message) < minMessageRunes {
		return nil, &domain.ValidationError{Message: "message is too short"}
	}
	if len(message) > maxMessageBytes {
		return nil, &domain.ValidationError{Message: "message is too long"}
	}

	rawEmail := strings.TrimSpace(input.ContactEmail)
	rawPhone := strings.TrimSpace(input.ContactPhone)
	if rawEmail == "" && rawPhone == "" {
		return nil, &domain.ValidationError{Message: "please provide an email or a Thai phone number to contact you back"}
	}

	var email, phone *string
	if rawEmail != "" {
		if !isValidEmail(rawEmail) {
			return nil, &domain.ValidationError{Message: "please enter a valid email address"}
		}
		v := rawEmail
		email = &v
	}
	if rawPhone != "" {
		if !isValidThaiPhone(rawPhone) {
			return nil, &domain.ValidationError{Message: "please enter a valid Thai phone number"}
		}
		v := normalizePhone(rawPhone)
		phone = &v
	}

	return &cleanSubmission{
		Name:         name,
		Message:      message,
		ContactEmail: email,
		ContactPhone: phone,
	}, nil
}
