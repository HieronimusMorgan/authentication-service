package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to []string, subject string, body string) {
	// SMTP server details
	smtpHost := "smtp.gmail.com" // Replace with your SMTP host (e.g., Mailtrap, Outlook)
	smtpPort := "587"

	// Sender credentials
	senderEmail := os.Getenv("USER_EMAIL")
	senderPassword := os.Getenv("PASS_EMAIL")

	// Construct message
	message := []byte(subject + "\n" + body)

	// Authenticate
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	// Send Email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, to, message)
	if err != nil {
		fmt.Println("❌ Failed to send email:", err)
		return
	}

	fmt.Println("✅ Email sent successfully!")
}
