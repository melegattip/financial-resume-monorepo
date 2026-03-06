package email

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// EmailService defines the contract for sending transactional emails.
type EmailService interface {
	SendPasswordReset(toEmail, resetLink string) error
	SendEmailVerification(toEmail, firstName, verificationLink string) error
}

// SMTPConfig holds SMTP connection settings.
type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

// ResendConfig holds Resend HTTP API settings.
type ResendConfig struct {
	APIKey string
	From   string
}

// SMTPEmailService sends emails via SMTP (Gmail / Google Workspace).
type SMTPEmailService struct {
	cfg    SMTPConfig
	logger zerolog.Logger
}

// ResendEmailService sends emails via Resend HTTP API (works on Render).
type ResendEmailService struct {
	cfg    ResendConfig
	client *http.Client
	logger zerolog.Logger
}

// NoOpEmailService logs reset links but does not send real emails.
// Used when SMTP is not configured (e.g. local dev without credentials).
type NoOpEmailService struct {
	logger zerolog.Logger
}

// NewService returns a ResendEmailService if a Resend API key is configured,
// an SMTPEmailService if SMTP credentials are set, or a NoOpEmailService otherwise.
func NewService(cfg SMTPConfig, logger zerolog.Logger) EmailService {
	return newService(cfg, ResendConfig{}, logger)
}

// NewServiceWithResend is like NewService but also accepts a Resend config.
// Resend takes priority over SMTP when an API key is present.
func NewServiceWithResend(smtp SMTPConfig, resend ResendConfig, logger zerolog.Logger) EmailService {
	return newService(smtp, resend, logger)
}

func newService(smtpCfg SMTPConfig, resendCfg ResendConfig, logger zerolog.Logger) EmailService {
	if resendCfg.APIKey != "" {
		logger.Info().Msg("email: using Resend HTTP API")
		return &ResendEmailService{
			cfg:    resendCfg,
			client: &http.Client{Timeout: 15 * time.Second},
			logger: logger,
		}
	}
	if smtpCfg.User != "" && smtpCfg.Password != "" {
		logger.Info().Msg("email: using SMTP")
		return &SMTPEmailService{cfg: smtpCfg, logger: logger}
	}
	logger.Warn().Msg("email: no credentials configured — using no-op (links will be logged)")
	return &NoOpEmailService{logger: logger}
}

const smtpDialTimeout = 15 * time.Second

// send dials SMTP and delivers a pre-built MIME message.
// Port 465 uses implicit TLS (SMTPS); any other port uses STARTTLS.
func (s *SMTPEmailService) send(toEmail, msg string) error {
	addr := s.cfg.Host + ":" + s.cfg.Port
	auth := smtp.PlainAuth("", s.cfg.User, s.cfg.Password, s.cfg.Host)
	tlsCfg := &tls.Config{ServerName: s.cfg.Host}

	var conn *smtp.Client
	if s.cfg.Port == "465" {
		// Implicit TLS (SMTPS) — TLS from the very first byte.
		tlsConn, err := tls.DialWithDialer(
			&net.Dialer{Timeout: smtpDialTimeout},
			"tcp", addr, tlsCfg,
		)
		if err != nil {
			return fmt.Errorf("smtp dial (tls): %w", err)
		}
		c, err := smtp.NewClient(tlsConn, s.cfg.Host)
		if err != nil {
			tlsConn.Close()
			return fmt.Errorf("smtp new client: %w", err)
		}
		conn = c
	} else {
		// STARTTLS (port 587 or other).
		netConn, err := net.DialTimeout("tcp", addr, smtpDialTimeout)
		if err != nil {
			return fmt.Errorf("smtp dial: %w", err)
		}
		c, err := smtp.NewClient(netConn, s.cfg.Host)
		if err != nil {
			netConn.Close()
			return fmt.Errorf("smtp new client: %w", err)
		}
		if err := c.StartTLS(tlsCfg); err != nil {
			c.Close()
			return fmt.Errorf("smtp starttls: %w", err)
		}
		conn = c
	}
	defer conn.Close()

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
	return nil
}

// SendPasswordReset sends a password-reset email.
func (s *SMTPEmailService) SendPasswordReset(toEmail, resetLink string) error {
	msg := buildMIMEMessage(s.cfg.From, toEmail, "Restablecer contraseña — Niloft", buildResetEmailHTML(resetLink))
	if err := s.send(toEmail, msg); err != nil {
		return err
	}
	s.logger.Info().Str("to", toEmail).Msg("password reset email sent")
	return nil
}

// SendEmailVerification sends an account verification email.
func (s *SMTPEmailService) SendEmailVerification(toEmail, firstName, verificationLink string) error {
	msg := buildMIMEMessage(s.cfg.From, toEmail, "Verificá tu cuenta — Niloft", buildVerificationEmailHTML(firstName, verificationLink))
	if err := s.send(toEmail, msg); err != nil {
		return err
	}
	s.logger.Info().Str("to", toEmail).Msg("verification email sent")
	return nil
}

// sendResend posts an email via the Resend HTTP API.
func (s *ResendEmailService) sendResend(toEmail, subject, htmlBody string) error {
	payload, err := json.Marshal(map[string]any{
		"from":    s.cfg.From,
		"to":      []string{toEmail},
		"subject": subject,
		"html":    htmlBody,
	})
	if err != nil {
		return fmt.Errorf("resend marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("resend request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("resend http: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("resend api returned status %d", resp.StatusCode)
	}
	return nil
}

// SendPasswordReset sends a password-reset email via Resend.
func (s *ResendEmailService) SendPasswordReset(toEmail, resetLink string) error {
	if err := s.sendResend(toEmail, "Restablecer contraseña — Niloft", buildResetEmailHTML(resetLink)); err != nil {
		return err
	}
	s.logger.Info().Str("to", toEmail).Msg("password reset email sent via Resend")
	return nil
}

// SendEmailVerification sends an account verification email via Resend.
func (s *ResendEmailService) SendEmailVerification(toEmail, firstName, verificationLink string) error {
	if err := s.sendResend(toEmail, "Verificá tu cuenta — Niloft", buildVerificationEmailHTML(firstName, verificationLink)); err != nil {
		return err
	}
	s.logger.Info().Str("to", toEmail).Msg("verification email sent via Resend")
	return nil
}

// SendPasswordReset logs the reset link (no-op mode).
func (s *NoOpEmailService) SendPasswordReset(toEmail, resetLink string) error {
	s.logger.Info().
		Str("to", toEmail).
		Str("reset_link", resetLink).
		Msg("[NO-OP EMAIL] password reset link — configure SMTP to send real emails")
	return nil
}

// SendEmailVerification logs the verification link (no-op mode).
func (s *NoOpEmailService) SendEmailVerification(toEmail, firstName, verificationLink string) error {
	s.logger.Info().
		Str("to", toEmail).
		Str("first_name", firstName).
		Str("verification_link", verificationLink).
		Msg("[NO-OP EMAIL] verification link — configure SMTP to send real emails")
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

func buildVerificationEmailHTML(firstName, verificationLink string) string {
	return `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; background: #f4f4f4; padding: 40px 0; margin: 0;">
  <div style="max-width: 480px; margin: 0 auto; background: #fff; border-radius: 12px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.08);">
    <div style="background: linear-gradient(135deg, #3b82f6, #6366f1); padding: 32px; text-align: center;">
      <img src="https://financial.niloft.com/logo-niloft.png" alt="Niloft" style="height: 48px; width: auto; display: block; margin: 0 auto;">
      <p style="color: rgba(255,255,255,0.8); margin: 10px 0 0; font-size: 13px;">Tu asistente financiero</p>
    </div>
    <div style="padding: 36px 32px;">
      <h2 style="color: #1f2937; margin-top: 0; font-size: 20px;">¡Bienvenido/a, ` + firstName + `!</h2>
      <p style="color: #6b7280; line-height: 1.6; margin-bottom: 28px;">
        Tu cuenta fue creada con éxito. Para activarla y empezar a gestionar tus finanzas,
        confirmá tu dirección de email haciendo clic en el botón.
      </p>
      <div style="text-align: center; margin-bottom: 28px;">
        <a href="` + verificationLink + `"
           style="background: #10b981; color: #fff; padding: 14px 36px; border-radius: 8px;
                  text-decoration: none; font-weight: 600; display: inline-block; font-size: 15px;">
          Verificar mi cuenta
        </a>
      </div>
      <p style="color: #9ca3af; font-size: 13px; text-align: center; line-height: 1.5;">
        Este enlace expira en <strong>24 horas</strong>.<br>
        Si no creaste esta cuenta, ignorá este correo.
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

func buildResetEmailHTML(resetLink string) string {
	return `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; background: #f4f4f4; padding: 40px 0; margin: 0;">
  <div style="max-width: 480px; margin: 0 auto; background: #fff; border-radius: 12px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.08);">
    <div style="background: linear-gradient(135deg, #3b82f6, #6366f1); padding: 32px; text-align: center;">
      <img src="https://financial.niloft.com/logo-niloft.png" alt="Niloft" style="height: 48px; width: auto; display: block; margin: 0 auto;">
      <p style="color: rgba(255,255,255,0.8); margin: 10px 0 0; font-size: 13px;">Tu asistente financiero</p>
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
