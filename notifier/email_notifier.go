package notifier

import (
	"bytes"
	"fmt"
	"net/smtp"

	"github.com/pyxsoft/notificator/event"
	"github.com/pyxsoft/notificator/preference"
	"github.com/pyxsoft/notificator/template"
)

type EmailNotifier struct {
	SMTPHost        string
	SMTPPort        int
	Username        string
	Password        string
	TemplateManager *template.TemplateManager
}

func NewEmailNotifier(smtpHost string, smtpPort int, username, password string) *EmailNotifier {
	return &EmailNotifier{
		SMTPHost: smtpHost,
		SMTPPort: smtpPort,
		Username: username,
		Password: password,
	}
}

func (n *EmailNotifier) GetChannelName() string {
	return "email"
}

func (n *EmailNotifier) SetTemplateManager(tmplManager *template.TemplateManager) {
	n.TemplateManager = tmplManager
}

func (n *EmailNotifier) Send(event *event.Event, userPrefs preference.UserPreferences) error {
	tmpl, err := n.TemplateManager.GetTemplate(userPrefs.PreferredLanguage, n.GetChannelName(), event.GetType())
	if err != nil {
		return err
	}

	dataMap, ok := event.GetData().(map[string]interface{})
	if !ok {
		return fmt.Errorf("event data is not a map")
	}

	recipient, ok := dataMap["Email"].(string)
	if !ok {
		return fmt.Errorf("email address not found in event data")
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, event.GetData()); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", n.Username, n.Password, n.SMTPHost)
	to := []string{recipient}
	msg := []byte(fmt.Sprintf("Subject: Notification\r\n\r\n%s", body.String()))

	addr := fmt.Sprintf("%s:%d", n.SMTPHost, n.SMTPPort)
	if err := smtp.SendMail(addr, auth, n.Username, to, msg); err != nil {
		return err
	}
	return nil
}
