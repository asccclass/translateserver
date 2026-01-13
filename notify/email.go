package notify

import (
	"fmt"
	"log"
	"net/smtp"
	"translateserver/config"
)

func SendEmail(to, status, content string) {
	if config.SmtpHost == "" || config.SmtpUser == "" || config.SmtpPass == "" {
		log.Printf("Email skipped (credentials missing). Mock Email to %s: Job %s. Content length: %d", to, status, len(content))
		return
	}

	// Basic SMTP implementation
	auth := smtp.PlainAuth("", config.SmtpUser, config.SmtpPass, config.SmtpHost)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: Translation %s\r\n\r\nStatus: %s\r\n\r\nOutput:\r\n%s", to, status, status, content))
	addr := fmt.Sprintf("%s:%s", config.SmtpHost, config.SmtpPort)

	err := smtp.SendMail(addr, auth, "noreply@translateserver.com", []string{to}, msg)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
	} else {
		log.Printf("Email sent to %s", to)
	}
}
