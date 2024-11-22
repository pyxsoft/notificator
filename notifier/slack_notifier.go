package notifier

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pyxsoft/notificator/event"
	"github.com/pyxsoft/notificator/preference"
	"github.com/pyxsoft/notificator/template"
)

type SlackNotifier struct {
	WebhookURL      string
	TemplateManager *template.TemplateManager
}

func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		WebhookURL: webhookURL,
	}
}

func (n *SlackNotifier) GetChannelName() string {
	return "slack"
}

func (n *SlackNotifier) SetTemplateManager(tmplManager *template.TemplateManager) {
	n.TemplateManager = tmplManager
}

func (n *SlackNotifier) Send(event *event.Event, userPrefs preference.UserPreferences) error {
	tmpl, err := n.TemplateManager.GetTemplate(userPrefs.PreferredLanguage, n.GetChannelName(), event.GetType())
	if err != nil {
		return err
	}

	var message bytes.Buffer
	if err := tmpl.Execute(&message, event.GetData()); err != nil {
		return err
	}

	payload := map[string]string{
		"text": message.String(),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(n.WebhookURL, "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
