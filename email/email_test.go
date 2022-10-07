package email

import (
	"testing"
)

// TestSendCompiles verifies that the email package compiles
func TestSendCompiles(t *testing.T) {
	smtpConfig := SMTPConfig{}
	mailer := NewMailer(smtpConfig)

	mail := Email{}
	mailer.Send(mail)
}
