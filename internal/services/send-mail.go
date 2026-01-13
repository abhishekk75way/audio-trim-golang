package services

import (
	"fmt"

	"authentication/backend/internal/utils"
)

func sendEmail(toEmail, token string) error {
	transporter, err := utils.NewEmailTransporter()
	if err != nil {
		return err
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", transporter.FrontendURL, token)

	body := fmt.Sprintf(
		"Click the link below to reset your password:\n\n%s\n\nIf you did not request this, ignore this email.",
		resetLink,
	)

	return transporter.SendPlain(
		toEmail,
		"Password Reset Request",
		body,
	)
}
