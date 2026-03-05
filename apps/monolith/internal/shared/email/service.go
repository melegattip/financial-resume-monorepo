package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/rs/zerolog"
)

// EmailService defines the contract for sending transactional emails.
type EmailService interface {
	SendPasswordReset(toEmail, resetLink string) error
}

// SMTPConfig holds SMTP connection settings.
type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

// SMTPEmailService sends emails via SMTP (Gmail / Google Workspace).
type SMTPEmailService struct {
	cfg    SMTPConfig
	logger zerolog.Logger
}

// NoOpEmailService logs reset links but does not send real emails.
// Used when SMTP is not configured (e.g. local dev without credentials).
type NoOpEmailService struct {
	logger zerolog.Logger
}

// NewService returns an SMTPEmailService if SMTP credentials are configured,
// or a NoOpEmailService otherwise.
func NewService(cfg SMTPConfig, logger zerolog.Logger) EmailService {
	if cfg.User == "" || cfg.Password == "" {
		logger.Warn().Msg("SMTP credentials not configured — using no-op email service (reset links will be logged)")
		return &NoOpEmailService{logger: logger}
	}
	return &SMTPEmailService{cfg: cfg, logger: logger}
}

// SendPasswordReset sends a password-reset email.
func (s *SMTPEmailService) SendPasswordReset(toEmail, resetLink string) error {
	subject := "Restablecer contraseña — Niloft"
	body := buildResetEmailHTML(resetLink)
	msg := buildMIMEMessage(s.cfg.From, toEmail, subject, body)

	addr := s.cfg.Host + ":" + s.cfg.Port
	auth := smtp.PlainAuth("", s.cfg.User, s.cfg.Password, s.cfg.Host)

	conn, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("smtp dial: %w", err)
	}
	defer conn.Close()

	if err := conn.StartTLS(&tls.Config{ServerName: s.cfg.Host}); err != nil {
		return fmt.Errorf("smtp starttls: %w", err)
	}
	if err := conn.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}
	if err := conn.Mail(s.cfg.From); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}
	if err := conn.Rcpt(toEmail); err != nil {
		return fmt.Errorf("smtp rcpt to: %w", err)
	}

	wc, err := conn.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	defer wc.Close()

	if _, err := fmt.Fprint(wc, msg); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}

	s.logger.Info().Str("to", toEmail).Msg("password reset email sent")
	return nil
}

// SendPasswordReset logs the reset link (no-op mode).
func (s *NoOpEmailService) SendPasswordReset(toEmail, resetLink string) error {
	s.logger.Info().
		Str("to", toEmail).
		Str("reset_link", resetLink).
		Msg("[NO-OP EMAIL] password reset link generated — configure SMTP to send real emails")
	return nil
}

func buildMIMEMessage(from, to, subject, htmlBody string) string {
	var sb strings.Builder
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	sb.WriteString("From: Niloft <" + from + ">\r\n")
	sb.WriteString("To: " + to + "\r\n")
	sb.WriteString("Subject: " + subject + "\r\n")
	sb.WriteString("\r\n")
	sb.WriteString(htmlBody)
	return sb.String()
}

func buildResetEmailHTML(resetLink string) string {
	return `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; background: #f4f4f4; padding: 40px 0; margin: 0;">
  <div style="max-width: 480px; margin: 0 auto; background: #fff; border-radius: 12px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.08);">
    <div style="background: linear-gradient(135deg, #3b82f6, #6366f1); padding: 32px; text-align: center;">
      <h1 style="color: #fff; margin: 0; font-size: 26px; letter-spacing: -0.5px;">Niloft</h1>
      <p style="color: rgba(255,255,255,0.8); margin: 6px 0 0; font-size: 13px;">Tu asistente financiero</p>
    </div>
    <div style="padding: 36px 32px;">
      <h2 style="color: #1f2937; margin-top: 0; font-size: 20px;">Restablecer contraseña</h2>
      <p style="color: #6b7280; line-height: 1.6; margin-bottom: 28px;">
        Recibimos una solicitud para restablecer la contraseña de tu cuenta.
        Hacé clic en el botón para crear una nueva contraseña.
      </p>
      <div style="text-align: center; margin-bottom: 28px;">
        <a href="` + resetLink + `"
           style="background: #3b82f6; color: #fff; padding: 14px 36px; border-radius: 8px;
                  text-decoration: none; font-weight: 600; display: inline-block; font-size: 15px;">
          Restablecer contraseña
        </a>
      </div>
      <p style="color: #9ca3af; font-size: 13px; text-align: center; line-height: 1.5;">
        Este enlace expira en <strong>1 hora</strong>.<br>
        Si no solicitaste este cambio, ignorá este correo.
      </p>
      <hr style="border: none; border-top: 1px solid #e5e7eb; margin: 24px 0;">
      <p style="color: #d1d5db; font-size: 12px; text-align: center; margin: 0;">
        © 2025 Niloft · Tu asistente financiero
      </p>
    </div>
  </div>
</body>
</html>`
}
