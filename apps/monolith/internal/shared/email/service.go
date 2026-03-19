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

// CoachingReportEmailData holds the data for the monthly coaching report email.
// It is defined here (not in the AI domain package) to avoid circular imports.
type CoachingReportEmailData struct {
	Month        string
	Sentiment    string // "positivo|neutral|desafiante"
	Summary      string
	Wins         []CoachingEmailPoint
	Improvements []CoachingEmailPoint
	Actions      []CoachingEmailAction
	BehaviorNote string
}

// CoachingEmailPoint is a single win or improvement item in the email.
type CoachingEmailPoint struct {
	Title       string
	Description string
}

// CoachingEmailAction is a concrete action item in the email.
type CoachingEmailAction struct {
	Title  string
	Detail string
}

// EmailService defines the contract for sending transactional emails.
type EmailService interface {
	SendPasswordReset(toEmail, resetLink string) error
	SendEmailVerification(toEmail, firstName, verificationLink string) error
	SendBudgetAlert(toEmail, firstName, categoryID, period, newStatus string, spentAmount, budgetLimit float64) error
	SendLoginNotification(toEmail, firstName, loginTime string) error
	SendMonthlyCoachingReport(toEmail, firstName string, data CoachingReportEmailData) error
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

// SendBudgetAlert sends a budget threshold alert via SMTP.
func (s *SMTPEmailService) SendBudgetAlert(toEmail, firstName, categoryID, period, newStatus string, spentAmount, budgetLimit float64) error {
	msg := buildMIMEMessage(s.cfg.From, toEmail, "Alerta de presupuesto — Niloft", buildBudgetAlertHTML(firstName, categoryID, period, newStatus, spentAmount, budgetLimit))
	if err := s.send(toEmail, msg); err != nil {
		return err
	}
	s.logger.Info().Str("to", toEmail).Msg("budget alert email sent")
	return nil
}

// SendLoginNotification sends a new login notification via SMTP.
func (s *SMTPEmailService) SendLoginNotification(toEmail, firstName, loginTime string) error {
	msg := buildMIMEMessage(s.cfg.From, toEmail, "Nuevo inicio de sesión — Niloft", buildLoginNotificationHTML(firstName, loginTime))
	if err := s.send(toEmail, msg); err != nil {
		return err
	}
	s.logger.Info().Str("to", toEmail).Msg("login notification email sent")
	return nil
}

// SendBudgetAlert sends a budget threshold alert via Resend.
func (s *ResendEmailService) SendBudgetAlert(toEmail, firstName, categoryID, period, newStatus string, spentAmount, budgetLimit float64) error {
	if err := s.sendResend(toEmail, "Alerta de presupuesto — Niloft", buildBudgetAlertHTML(firstName, categoryID, period, newStatus, spentAmount, budgetLimit)); err != nil {
		return err
	}
	s.logger.Info().Str("to", toEmail).Msg("budget alert email sent via Resend")
	return nil
}

// SendLoginNotification sends a new login notification via Resend.
func (s *ResendEmailService) SendLoginNotification(toEmail, firstName, loginTime string) error {
	if err := s.sendResend(toEmail, "Nuevo inicio de sesión — Niloft", buildLoginNotificationHTML(firstName, loginTime)); err != nil {
		return err
	}
	s.logger.Info().Str("to", toEmail).Msg("login notification email sent via Resend")
	return nil
}

// SendBudgetAlert logs a budget alert (no-op mode).
func (s *NoOpEmailService) SendBudgetAlert(toEmail, firstName, categoryID, period, newStatus string, spentAmount, budgetLimit float64) error {
	s.logger.Info().
		Str("to", toEmail).
		Str("category_id", categoryID).
		Str("new_status", newStatus).
		Msg("[NO-OP EMAIL] budget alert — configure email to send real emails")
	return nil
}

// SendLoginNotification logs a login notification (no-op mode).
func (s *NoOpEmailService) SendLoginNotification(toEmail, firstName, loginTime string) error {
	s.logger.Info().
		Str("to", toEmail).
		Str("login_time", loginTime).
		Msg("[NO-OP EMAIL] login notification — configure email to send real emails")
	return nil
}

// SendMonthlyCoachingReport sends the monthly coaching report via SMTP.
func (s *SMTPEmailService) SendMonthlyCoachingReport(toEmail, firstName string, data CoachingReportEmailData) error {
	subject := "Tu resumen de " + data.Month + " — Niloft"
	msg := buildMIMEMessage(s.cfg.From, toEmail, subject, buildMonthlyCoachingHTML(firstName, data))
	if err := s.send(toEmail, msg); err != nil {
		return err
	}
	s.logger.Info().Str("to", toEmail).Str("month", data.Month).Msg("monthly coaching report email sent")
	return nil
}

// SendMonthlyCoachingReport sends the monthly coaching report via Resend.
func (s *ResendEmailService) SendMonthlyCoachingReport(toEmail, firstName string, data CoachingReportEmailData) error {
	subject := "Tu resumen de " + data.Month + " — Niloft"
	if err := s.sendResend(toEmail, subject, buildMonthlyCoachingHTML(firstName, data)); err != nil {
		return err
	}
	s.logger.Info().Str("to", toEmail).Str("month", data.Month).Msg("monthly coaching report email sent via Resend")
	return nil
}

// SendMonthlyCoachingReport logs the coaching report (no-op mode).
func (s *NoOpEmailService) SendMonthlyCoachingReport(toEmail, firstName string, data CoachingReportEmailData) error {
	s.logger.Info().
		Str("to", toEmail).
		Str("month", data.Month).
		Str("sentiment", data.Sentiment).
		Msg("[NO-OP EMAIL] monthly coaching report — configure email to send real emails")
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

func buildBudgetAlertHTML(firstName, categoryID, period, newStatus string, spentAmount, budgetLimit float64) string {
	spentPct := 0.0
	if budgetLimit > 0 {
		spentPct = (spentAmount / budgetLimit) * 100
	}
	statusLabel := "en alerta"
	statusColor := "#f59e0b"
	if newStatus == "exceeded" {
		statusLabel = "excedido"
		statusColor = "#ef4444"
	}
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; background: #f4f4f4; padding: 40px 0; margin: 0;">
  <div style="max-width: 480px; margin: 0 auto; background: #fff; border-radius: 12px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.08);">
    <div style="background: linear-gradient(135deg, #3b82f6, #6366f1); padding: 32px; text-align: center;">
      <img src="https://financial.niloft.com/logo-niloft.png" alt="Niloft" style="height: 48px; width: auto; display: block; margin: 0 auto;">
      <p style="color: rgba(255,255,255,0.8); margin: 10px 0 0; font-size: 13px;">Tu asistente financiero</p>
    </div>
    <div style="padding: 36px 32px;">
      <h2 style="color: #1f2937; margin-top: 0; font-size: 20px;">Alerta de presupuesto, %s</h2>
      <p style="color: #6b7280; line-height: 1.6; margin-bottom: 20px;">
        Tu presupuesto para el período <strong>%s</strong> está <span style="color:%s;font-weight:600;">%s</span>.
      </p>
      <div style="background:#f9fafb;border-radius:8px;padding:20px;margin-bottom:24px;">
        <p style="margin:0 0 8px;color:#374151;font-size:14px;"><strong>Categoría:</strong> %s</p>
        <p style="margin:0 0 8px;color:#374151;font-size:14px;"><strong>Gastado:</strong> $%.2f de $%.2f</p>
        <p style="margin:0;color:%s;font-size:14px;font-weight:600;"><strong>Porcentaje:</strong> %.1f%%</p>
      </div>
      <hr style="border: none; border-top: 1px solid #e5e7eb; margin: 24px 0;">
      <p style="color: #d1d5db; font-size: 12px; text-align: center; margin: 0;">
        © 2025 Niloft · Tu asistente financiero
      </p>
    </div>
  </div>
</body>
</html>`, firstName, period, statusColor, statusLabel, categoryID, spentAmount, budgetLimit, statusColor, spentPct)
}

func buildLoginNotificationHTML(firstName, loginTime string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; background: #f4f4f4; padding: 40px 0; margin: 0;">
  <div style="max-width: 480px; margin: 0 auto; background: #fff; border-radius: 12px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.08);">
    <div style="background: linear-gradient(135deg, #3b82f6, #6366f1); padding: 32px; text-align: center;">
      <img src="https://financial.niloft.com/logo-niloft.png" alt="Niloft" style="height: 48px; width: auto; display: block; margin: 0 auto;">
      <p style="color: rgba(255,255,255,0.8); margin: 10px 0 0; font-size: 13px;">Tu asistente financiero</p>
    </div>
    <div style="padding: 36px 32px;">
      <h2 style="color: #1f2937; margin-top: 0; font-size: 20px;">Nuevo inicio de sesión, %s</h2>
      <p style="color: #6b7280; line-height: 1.6; margin-bottom: 20px;">
        Se detectó un nuevo inicio de sesión en tu cuenta el <strong>%s</strong>.
      </p>
      <p style="color: #6b7280; line-height: 1.6; font-size: 14px;">
        Si no fuiste vos, cambiá tu contraseña inmediatamente desde la sección de seguridad.
      </p>
      <hr style="border: none; border-top: 1px solid #e5e7eb; margin: 24px 0;">
      <p style="color: #d1d5db; font-size: 12px; text-align: center; margin: 0;">
        © 2025 Niloft · Tu asistente financiero
      </p>
    </div>
  </div>
</body>
</html>`, firstName, loginTime)
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

func buildMonthlyCoachingHTML(firstName string, data CoachingReportEmailData) string {
	var sb strings.Builder

	// Sentiment badge config
	sentimentLabel := "Mes Neutral"
	sentimentBg := "#fef3c7"
	sentimentColor := "#92400e"
	sentimentEmoji := "⚖️"
	switch data.Sentiment {
	case "positivo":
		sentimentLabel = "Mes Positivo"
		sentimentBg = "#d1fae5"
		sentimentColor = "#065f46"
		sentimentEmoji = "✨"
	case "desafiante":
		sentimentLabel = "Mes Desafiante"
		sentimentBg = "#fee2e2"
		sentimentColor = "#991b1b"
		sentimentEmoji = "💪"
	}

	sb.WriteString(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; background: #f4f4f4; padding: 40px 0; margin: 0;">
  <div style="max-width: 480px; margin: 0 auto; background: #fff; border-radius: 12px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.08);">
    <div style="background: linear-gradient(135deg, #3b82f6, #6366f1); padding: 32px; text-align: center;">
      <img src="https://financial.niloft.com/logo-niloft.png" alt="Niloft" style="height: 48px; width: auto; display: block; margin: 0 auto;">
      <p style="color: rgba(255,255,255,0.8); margin: 10px 0 4px; font-size: 12px; text-transform: uppercase; letter-spacing: 1px;">Reporte Mensual</p>
      <p style="color: #fff; margin: 0; font-size: 22px; font-weight: 700;">` + data.Month + `</p>
    </div>
    <div style="padding: 28px 32px 0;">`)

	sb.WriteString(fmt.Sprintf(
		`<div style="display:inline-block;background:%s;color:%s;font-size:12px;font-weight:600;padding:4px 12px;border-radius:20px;margin-bottom:16px;">%s %s</div>`,
		sentimentBg, sentimentColor, sentimentEmoji, sentimentLabel))

	if firstName != "" {
		sb.WriteString(fmt.Sprintf(
			`<p style="color:#1f2937;font-size:15px;font-weight:600;margin:0 0 8px;">Hola, %s</p>`, firstName))
	}
	sb.WriteString(fmt.Sprintf(
		`<p style="color:#374151;font-size:14px;line-height:1.7;margin:0 0 24px;">%s</p></div>`, data.Summary))

	// Wins
	if len(data.Wins) > 0 {
		sb.WriteString(`<div style="padding:0 32px 20px;"><p style="font-size:13px;font-weight:700;color:#4f46e5;text-transform:uppercase;letter-spacing:0.8px;margin:0 0 12px;">✅ Lo que hiciste bien</p>`)
		for _, w := range data.Wins {
			sb.WriteString(fmt.Sprintf(
				`<div style="background:#f0fdf4;border-left:3px solid #10b981;border-radius:0 8px 8px 0;padding:10px 14px;margin-bottom:8px;"><p style="font-size:13px;font-weight:600;color:#065f46;margin:0 0 2px;">%s</p><p style="font-size:13px;color:#374151;line-height:1.5;margin:0;">%s</p></div>`,
				w.Title, w.Description))
		}
		sb.WriteString(`</div>`)
	}

	// Improvements
	if len(data.Improvements) > 0 {
		sb.WriteString(`<div style="padding:0 32px 20px;"><p style="font-size:13px;font-weight:700;color:#4f46e5;text-transform:uppercase;letter-spacing:0.8px;margin:0 0 12px;">🎯 Para mejorar</p>`)
		for _, imp := range data.Improvements {
			sb.WriteString(fmt.Sprintf(
				`<div style="background:#fffbeb;border-left:3px solid #f59e0b;border-radius:0 8px 8px 0;padding:10px 14px;margin-bottom:8px;"><p style="font-size:13px;font-weight:600;color:#92400e;margin:0 0 2px;">%s</p><p style="font-size:13px;color:#374151;line-height:1.5;margin:0;">%s</p></div>`,
				imp.Title, imp.Description))
		}
		sb.WriteString(`</div>`)
	}

	// Actions
	if len(data.Actions) > 0 {
		sb.WriteString(`<div style="padding:0 32px 20px;"><p style="font-size:13px;font-weight:700;color:#4f46e5;text-transform:uppercase;letter-spacing:0.8px;margin:0 0 12px;">🚀 Acciones para este mes</p>`)
		for _, a := range data.Actions {
			sb.WriteString(fmt.Sprintf(
				`<div style="border:1px solid #e5e7eb;border-radius:10px;padding:12px 14px;margin-bottom:8px;"><p style="font-size:13px;font-weight:600;color:#1f2937;margin:0 0 2px;">%s</p><p style="font-size:13px;color:#6b7280;line-height:1.5;margin:0;">%s</p></div>`,
				a.Title, a.Detail))
		}
		sb.WriteString(`</div>`)
	}

	// Behavior note
	if data.BehaviorNote != "" {
		sb.WriteString(fmt.Sprintf(
			`<div style="padding:0 32px 24px;"><div style="background:#eef2ff;border-radius:10px;padding:14px 16px;"><p style="font-size:12px;font-weight:600;color:#4f46e5;margin:0 0 4px;">📊 Tu patrón financiero</p><p style="font-size:13px;color:#374151;line-height:1.5;margin:0;">%s</p></div></div>`,
			data.BehaviorNote))
	}

	sb.WriteString(`<div style="padding:0 32px 32px;text-align:center;"><a href="https://financial.niloft.com/insights" style="display:inline-block;background:linear-gradient(135deg,#4f46e5,#6366f1);color:#fff;font-size:14px;font-weight:600;padding:12px 28px;border-radius:8px;text-decoration:none;margin-bottom:16px;">Ver reporte completo →</a><p style="color:#9ca3af;font-size:11px;margin:0;">© 2025 Niloft · Tu asistente financiero</p></div></div></body></html>`)

	return sb.String()
}
