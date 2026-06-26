package mailer

import (
	"fmt"
	"log/slog"
	"net/smtp"

	"notification-service/config"
)

func SendReceiptEmail(cfg *config.Config, to string, orderID string) error {
	if cfg.SMTPHost == "" || cfg.SMTPPort == "" {
		slog.Warn("SMTP configuration is missing. Simulating email send instead.", "to", to, "order_id", orderID)
		fmt.Printf("\n--- [SIMULATED EMAIL] ---\nTo: %s\nSubject: Pembayaran Berhasil - Pesanan #%s\nBody: Terima kasih! Pembayaran untuk pesanan Anda telah berhasil.\n-------------------------\n\n", to, orderID)
		return nil
	}

	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)

	from := "noreply@fooddelivery.com"
	subject := "Pembayaran Berhasil - Pesanan #" + orderID
	
	// Create HTML body
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Terima kasih atas pesanan Anda!</h2>
			<p>Pembayaran untuk pesanan <b>#%s</b> telah berhasil kami terima.</p>
			<p>Restoran sedang menyiapkan pesanan Anda dan akan segera diantarkan.</p>
			<br/>
			<p>Salam hangat,</p>
			<p>Tim Food Delivery</p>
		</body>
		</html>
	`, orderID)

	msg := []byte("To: " + to + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
		body + "\r\n")

	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)
	
	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		slog.Error("Failed to send email via SMTP", "error", err, "to", to)
		return err
	}

	slog.Info("Successfully sent email receipt", "to", to, "order_id", orderID)
	return nil
}
