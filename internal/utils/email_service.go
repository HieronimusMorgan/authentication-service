package utils

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

func SendEmail(to []string, subject string, body string) {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	senderEmail := os.Getenv("USER_EMAIL")
	senderPassword := os.Getenv("PASS_EMAIL") // Use App Password, not real password

	message := []byte("From: " + senderEmail + "\r\n" +
		"To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		body)

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, to, message)
	if err != nil {
		fmt.Println("❌ Failed to send email:", err)
		return
	}

	fmt.Println("✅ Email sent successfully!")
}
