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

type CustomPostNotifier struct {
	URL string
}

func NewCustomPostNotifier(url string) *CustomPostNotifier {
	return &CustomPostNotifier{
		URL: url,
	}
}

func (n *CustomPostNotifier) GetChannelName() string {
	return "custom_post"
}

func (n *CustomPostNotifier) SetTemplateManager(_ *template.TemplateManager) {
	// No template manager needed for this notifier
}

func (n *CustomPostNotifier) Send(event *event.Event, userPrefs preference.UserPreferences) error {
	payload := map[string]interface{}{
		"event_type": event.GetType(),
		"data":       event.GetData(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(n.URL, "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("custom URL returned status code %d", resp.StatusCode)
	}

	return nil
}
