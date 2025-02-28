package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func SendTopupNotification(email, userName string, amount, newBalance float64, paymentID string) error {
	apiKey := os.Getenv("MAILERSEND_API_KEY")
	fromEmail := os.Getenv("FROM_EMAIL")
	fromName := os.Getenv("FROM_NAME")

	if fromEmail == "" {
		fromEmail = "noreply@trial-yzkq340vnnkld796.mlsender.net"
	}
	if fromName == "" {
		fromName = "Game Rental Service"
	}
	if apiKey == "" {
		return fmt.Errorf("MAILERSEND_API_KEY is not set")
	}

	htmlContent := fmt.Sprintf(`
		<html>
		<body>
			<h1>Topup Successful!</h1>
			<p>Dear %s,</p>
			<p>Your account has been successfully charged with <strong>$%.2f</strong>.</p>
			<p>Your new balance is: <strong>$%.2f</strong>.</p>
			<p>Payment ID: <strong>%s</strong></p>
			<p>Thank you for using our service!</p>
			<p>Regards,<br>Video Game Rental Team</p>
		</body>
		</html>
	`, userName, amount, newBalance, paymentID)

	textContent := fmt.Sprintf(
		"Topup Successful!\n\nDear %s,\n\nYour account has been successfully charged with $%.2f.\nYour new balance is: $%.2f.\nPayment ID: %s\n\nThank you for using our service!\n\nRegards,\nVideo Game Rental Team",
		userName, amount, newBalance, paymentID)

	type EmailAddress struct {
		Email string `json:"email"`
		Name  string `json:"name,omitempty"`
	}

	type PersonalizationData struct {
		Email string                 `json:"email"`
		Data  map[string]interface{} `json:"data,omitempty"`
	}

	type EmailRequest struct {
		From        EmailAddress          `json:"from"`
		To          []EmailAddress        `json:"to"`
		Subject     string                `json:"subject"`
		HTML        string                `json:"html"`
		Text        string                `json:"text"`
		Tags        []string              `json:"tags,omitempty"`
		Personalize []PersonalizationData `json:"personalization,omitempty"`
	}

	payload := EmailRequest{
		From: EmailAddress{
			Email: fromEmail,
			Name:  fromName,
		},
		To: []EmailAddress{
			{
				Email: "zuyatna@gmail.com",
				Name:  userName,
			},
		},
		Subject: "Topup Successful",
		HTML:    htmlContent,
		Text:    textContent,
		Tags:    []string{"topup", "notification"},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.mailersend.com/v1/email", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("MailerSend API returned error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func SendBookingNotification(email, userName string, status string, booking_id int, total_pay float64) error {
	apiKey := os.Getenv("MAILERSEND_API_KEY")
	fromEmail := os.Getenv("FROM_EMAIL")
	fromName := os.Getenv("FROM_NAME")

	if fromEmail == "" {
		fromEmail = "noreply@trial-yzkq340vnnkld796.mlsender.net"
	}
	if fromName == "" {
		fromName = "Game Rental Service"
	}
	if apiKey == "" {
		return fmt.Errorf("MAILERSEND_API_KEY is not set")
	}

	htmlContent := fmt.Sprintf(`
		<html>
		<body>
			<h1>Booking %s</h1>
			<p>Dear %s,</p>
			<p>Your game rental booking has been <strong>%s</strong>.</p>
			<p>Booking ID: <strong>%d</strong></p>
			<p>Total Payment: <strong>$%.2f</strong></p>
			<p>Thank you for using our service!</p>
			<p>Regards,<br>Video Game Rental Team</p>
		</body>
		</html>
	`, status, userName, status, booking_id, total_pay)

	textContent := fmt.Sprintf(
		"Booking %s\n\nDear %s,\n\nYour game rental booking has been %s.\nAmount: $%.2f\n\nThank you for using our service!\n\nRegards,\nVideo Game Rental Team",
		status, userName, status, total_pay)

	type EmailAddress struct {
		Email string `json:"email"`
		Name  string `json:"name,omitempty"`
	}

	type PersonalizationData struct {
		Email string                 `json:"email"`
		Data  map[string]interface{} `json:"data,omitempty"`
	}

	type EmailRequest struct {
		From        EmailAddress          `json:"from"`
		To          []EmailAddress        `json:"to"`
		Subject     string                `json:"subject"`
		HTML        string                `json:"html"`
		Text        string                `json:"text"`
		Tags        []string              `json:"tags,omitempty"`
		Personalize []PersonalizationData `json:"personalization,omitempty"`
	}

	payload := EmailRequest{
		From: EmailAddress{
			Email: fromEmail,
			Name:  fromName,
		},
		To: []EmailAddress{
			{
				Email: "zuyatna@gmail.com",
				Name:  userName,
			},
		},
		Subject: fmt.Sprintf("Game Rental Booking %s", status),
		HTML:    htmlContent,
		Text:    textContent,
		Tags:    []string{"booking", "notification"},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.mailersend.com/v1/email", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("MailerSend API returned error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func SendTransactionNotification(email, userName string, transferID int, amount float64, receiverID uuid.UUID, balance float64) error {
	apiKey := os.Getenv("MAILERSEND_API_KEY")
	fromEmail := os.Getenv("FROM_EMAIL")
	fromName := os.Getenv("FROM_NAME")

	if fromEmail == "" {
		fromEmail = "noreply@trial-yzkq340vnnkld796.mlsender.net"
	}
	if fromName == "" {
		fromName = "Game Rental Service"
	}
	if apiKey == "" {
		return fmt.Errorf("MAILERSEND_API_KEY is not set")
	}

	htmlContent := fmt.Sprintf(`
		<html>
		<body>
			<h1>Fund Transfer</h1>
			<p>Dear %s,</p>
			<p>Your wallet has been transfered.</p>
			<p>Transfer ID: <strong>%d</strong></p>
			<p>Amount Transferred: <strong>$%.2f</strong></p>
			<p>To: <strong>%s</strong></p>
			<p>Your Balance: <strong>%.2f</strong></p>
			<p>Thank you for using our service!</p>
			<p>Regards,<br>Video Game Rental Team</p>
		</body>
		</html>
	`, userName, transferID, amount, receiverID, balance)

	textContent := fmt.Sprintf(
		"Fund Transfer\n\nDear %s,\n\nYour wallet has been transfered.\nTransfer ID: %d\nAmount Transferred: $%.2f\nto: %s\n\nThank you for using our service!\n\nRegards,\nVideo Game Rental Team",
		userName, transferID, amount, receiverID)

	type EmailAddress struct {
		Email string `json:"email"`
		Name  string `json:"name,omitempty"`
	}

	type PersonalizationData struct {
		Email string                 `json:"email"`
		Data  map[string]interface{} `json:"data,omitempty"`
	}

	type EmailRequest struct {
		From        EmailAddress          `json:"from"`
		To          []EmailAddress        `json:"to"`
		Subject     string                `json:"subject"`
		HTML        string                `json:"html"`
		Text        string                `json:"text"`
		Tags        []string              `json:"tags,omitempty"`
		Personalize []PersonalizationData `json:"personalization,omitempty"`
	}

	payload := EmailRequest{
		From: EmailAddress{
			Email: fromEmail,
			Name:  fromName,
		},
		To: []EmailAddress{
			{
				Email: "zuyatna@gmail.com",
				Name:  userName,
			},
		},
		Subject: "Fund Transfer",
		HTML:    htmlContent,
		Text:    textContent,
		Tags:    []string{"transfer", "notification"},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.mailersend.com/v1/email", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("MailerSend API returned error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
