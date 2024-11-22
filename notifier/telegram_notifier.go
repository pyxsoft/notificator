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

type TelegramNotifier struct {
	TemplateManager *template.TemplateManager
	BotToken        string
	ChatID          string
}

func NewTelegramNotifier(botToken, chatID string) *TelegramNotifier {
	return &TelegramNotifier{
		BotToken: botToken,
		ChatID:   chatID,
	}
}

func (n *TelegramNotifier) GetChannelName() string {
	return "telegram"
}

func (n *TelegramNotifier) SetTemplateManager(tmplManager *template.TemplateManager) {
	n.TemplateManager = tmplManager
}

func (n *TelegramNotifier) Send(event *event.Event, userPrefs preference.UserPreferences) error {
	tmpl, err := n.TemplateManager.GetTemplate(userPrefs.PreferredLanguage, n.GetChannelName(), event.GetType())
	if err != nil {
		return err
	}

	var message bytes.Buffer
	if err := tmpl.Execute(&message, event.GetData()); err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.BotToken)

	payload := map[string]string{
		"chat_id": n.ChatID,
		"text":    message.String(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Telegram API returned status code %d", resp.StatusCode)
	}

	return nil
}
