package email

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/smtp"
	"os"
	"strings"
)

type Sender struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
	Enabled  bool
}

func NewSender() *Sender {
	return &Sender{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		User:     os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
		Enabled:  os.Getenv("SMTP_ENABLED") == "true",
	}
}

func (s *Sender) IsConfigured() bool {
	return s.Enabled && s.Host != "" && s.From != ""
}

func (s *Sender) Send(to, subject, htmlBody string) error {
	if !s.IsConfigured() {
		slog.Warn("email not configured, skipping", "to", to, "subject", subject)
		return nil
	}

	// Sanitize headers to prevent injection
	sanitize := func(s string) string {
		return strings.NewReplacer("\r", "", "\n", "").Replace(s)
	}
	headers := map[string]string{
		"From":         s.From,
		"To":           sanitize(to),
		"Subject":      sanitize(subject),
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)

	var auth smtp.Auth
	if s.User != "" {
		auth = smtp.PlainAuth("", s.User, s.Password, s.Host)
	}

	tlsConfig := &tls.Config{ServerName: s.Host}

	if s.Port == "465" {
		// SSL
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("tls dial: %w", err)
		}
		client, err := smtp.NewClient(conn, s.Host)
		if err != nil {
			return fmt.Errorf("smtp client: %w", err)
		}
		defer client.Close()
		if auth != nil {
			if err := client.Auth(auth); err != nil {
				return fmt.Errorf("smtp auth: %w", err)
			}
		}
		if err := client.Mail(s.From); err != nil {
			return err
		}
		if err := client.Rcpt(to); err != nil {
			return err
		}
		w, err := client.Data()
		if err != nil {
			return err
		}
		w.Write([]byte(msg.String()))
		w.Close()
		return client.Quit()
	}

	// STARTTLS (port 587) or plain (port 25)
	return smtp.SendMail(addr, auth, s.From, []string{to}, []byte(msg.String()))
}

// SendAsync sends email in a goroutine so it doesn't block the HTTP response
func (s *Sender) SendAsync(to, subject, htmlBody string) {
	go func() {
		if err := s.Send(to, subject, htmlBody); err != nil {
			slog.Error("failed to send email", "to", to, "error", err)
		} else {
			slog.Info("email sent", "to", to, "subject", subject)
		}
	}()
}

func (s *Sender) SendCVApproved(to, cvTitle, lang string) {
	subject := "CV Approved - تمت الموافقة على السيرة الذاتية"
	body := fmt.Sprintf(`<div dir="rtl" style="font-family: Arial, sans-serif; padding: 20px;">
		<h2 style="color: #28a745;">✅ تمت الموافقة على سيرتك الذاتية</h2>
		<p>تمت الموافقة على السيرة الذاتية: <strong>%s</strong></p>
		<p>يمكنك الآن مشاركتها أو تحميلها بصيغة PDF.</p>
		<hr/>
		<h2 style="color: #28a745;">✅ Your CV has been approved</h2>
		<p>Your CV <strong>%s</strong> has been approved.</p>
		<p>You can now share it or download it as PDF.</p>
	</div>`, cvTitle, cvTitle)
	s.SendAsync(to, subject, body)
}

func (s *Sender) SendCVRejected(to, cvTitle, reason, lang string) {
	subject := "CV Rejected - تم رفض السيرة الذاتية"
	body := fmt.Sprintf(`<div dir="rtl" style="font-family: Arial, sans-serif; padding: 20px;">
		<h2 style="color: #dc3545;">❌ تم رفض سيرتك الذاتية</h2>
		<p>السيرة الذاتية: <strong>%s</strong></p>
		<p>السبب: %s</p>
		<hr/>
		<h2 style="color: #dc3545;">❌ Your CV has been rejected</h2>
		<p>CV: <strong>%s</strong></p>
		<p>Reason: %s</p>
	</div>`, cvTitle, reason, cvTitle, reason)
	s.SendAsync(to, subject, body)
}

func (s *Sender) SendCVRevision(to, cvTitle, notes, lang string) {
	subject := "CV Revision Required - مطلوب تعديل السيرة الذاتية"
	body := fmt.Sprintf(`<div dir="rtl" style="font-family: Arial, sans-serif; padding: 20px;">
		<h2 style="color: #ffc107;">📝 مطلوب تعديل سيرتك الذاتية</h2>
		<p>السيرة الذاتية: <strong>%s</strong></p>
		<p>ملاحظات: %s</p>
		<hr/>
		<h2 style="color: #ffc107;">📝 Your CV needs revision</h2>
		<p>CV: <strong>%s</strong></p>
		<p>Notes: %s</p>
	</div>`, cvTitle, notes, cvTitle, notes)
	s.SendAsync(to, subject, body)
}
