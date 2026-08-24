package usecase

import "testing"

func TestIsValidThaiPhone(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{"mobile 08", "0812345678", true},
		{"mobile 06", "0612345678", true},
		{"mobile 09", "0912345678", true},
		{"landline bangkok", "021234567", true},
		{"landline with dashes", "02-123-4567", true},
		{"mobile with spaces", "081 234 5678", true},
		{"mobile with parens", "(081) 234-5678", true},
		{"invalid prefix 00", "0012345678", false},
		{"invalid prefix 01", "0112345678", false},
		{"too short", "081234", false},
		{"too long", "081234567890", false},
		{"missing leading zero", "812345678", false},
		{"contains letters", "081abc5678", false},
		{"empty", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := isValidThaiPhone(tc.input)
			if got != tc.want {
				t.Errorf("isValidThaiPhone(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{"bare address", "user@example.com", true},
		{"display name form rejected", "Attacker <a@x.com>", false},
		{"missing at", "not-an-email", false},
		{"empty", "", false},
		{"trailing space tolerated", " user@example.com ", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := isValidEmail(tc.input)
			if got != tc.want {
				t.Errorf("isValidEmail(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestSanitizeHeaderValue(t *testing.T) {
	input := "สวัสดี\r\nBcc: spam-list@evil.com\r\nSubject: FREE MONEY"
	got := sanitizeHeaderValue(input)
	if got != "สวัสดีBcc: spam-list@evil.comSubject: FREE MONEY" {
		t.Errorf("sanitizeHeaderValue did not strip CRLF, got: %q", got)
	}
}

func TestStripControlChars(t *testing.T) {
	input := "hello\x00world"
	got := stripControlChars(input)
	if got != "helloworld" {
		t.Errorf("stripControlChars did not strip NUL byte, got: %q", got)
	}
}
