package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pyxsoft/notificator/event"
	"github.com/pyxsoft/notificator/preference"
	"github.com/pyxsoft/notificator/template"
)

type MailchimpNotifier struct {
	TemplateManager *template.TemplateManager
	APIKey          string
}

func NewMailchimpNotifier(apiKey string) *MailchimpNotifier {
	return &MailchimpNotifier{
		APIKey: apiKey,
	}
}

func (n *MailchimpNotifier) GetChannelName() string {
	return "mailchimp"
}

func (n *MailchimpNotifier) SetTemplateManager(tmplManager *template.TemplateManager) {
	n.TemplateManager = tmplManager
}

func (n *MailchimpNotifier) Send(event *event.Event, userPrefs preference.UserPreferences) error {
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

	var content bytes.Buffer
	if err := tmpl.Execute(&content, event.GetData()); err != nil {
		return err
	}

	payload := map[string]interface{}{
		"key": n.APIKey,
		"message": map[string]interface{}{
			"html":       content.String(),
			"subject":    "Notification",
			"from_email": "no-reply@example.com",
			"to": []map[string]interface{}{
				{
					"email": recipient,
					"type":  "to",
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post("https://mandrillapp.com/api/1.0/messages/send.json", "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Mailchimp API returned status code %d", resp.StatusCode)
	}

	return nil
}
