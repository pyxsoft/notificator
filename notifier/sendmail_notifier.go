package notifier

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pyxsoft/notificator/event"
	"github.com/pyxsoft/notificator/preference"
	"github.com/pyxsoft/notificator/template"
)

// SendmailNotifier sends emails using the local sendmail command.
type SendmailNotifier struct {
	FromAddress     string
	TemplateManager *template.TemplateManager
}

// NewSendmailNotifier creates a new SendmailNotifier.
func NewSendmailNotifier(fromAddress string) *SendmailNotifier {
	return &SendmailNotifier{
		FromAddress: fromAddress,
	}
}

// GetChannelName returns the name of the channel.
func (n *SendmailNotifier) GetChannelName() string {
	return "sendmail"
}

// SetTemplateManager sets the TemplateManager for the notifier.
func (n *SendmailNotifier) SetTemplateManager(tmplManager *template.TemplateManager) {
	n.TemplateManager = tmplManager
}

// Send sends a notification using the sendmail command.
func (n *SendmailNotifier) Send(event *event.Event, userPrefs preference.UserPreferences) error {
	// Get the template for the event type and user preferences
	tmpl, err := n.TemplateManager.GetTemplate(userPrefs.PreferredLanguage, n.GetChannelName(), event.GetType())
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	// Ensure event data is a map
	dataMap, ok := event.GetData().(map[string]interface{})
	if !ok {
		return fmt.Errorf("event data is not a map")
	}

	// Extract recipient from event data
	recipient, ok := dataMap["Email"].(string)
	if !ok || strings.TrimSpace(recipient) == "" {
		return fmt.Errorf("email address not found or empty in event data")
	}

	// Execute the template with event data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, dataMap); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Prepare the email content
	emailContent := fmt.Sprintf("From: %s\nTo: %s\nSubject: Notification\n\n%s",
		n.FromAddress, recipient, body.String())

	// Execute the sendmail command
	cmd := exec.Command("sendmail", "-t", "-oi")
	cmd.Stdin = strings.NewReader(emailContent)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute sendmail: %w", err)
	}

	return nil
}
