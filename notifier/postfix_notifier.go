package notifier

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/pyxsoft/notificator/event"
	"github.com/pyxsoft/notificator/preference"
	"github.com/pyxsoft/notificator/template"
)

type PostfixNotifier struct {
	TemplateManager *template.TemplateManager
	FromAddress     string
}

func NewPostfixNotifier(fromAddress string) *PostfixNotifier {
	return &PostfixNotifier{
		FromAddress: fromAddress,
	}
}

func (n *PostfixNotifier) GetChannelName() string {
	return "postfix"
}

func (n *PostfixNotifier) SetTemplateManager(tmplManager *template.TemplateManager) {
	n.TemplateManager = tmplManager
}

func (n *PostfixNotifier) Send(event *event.Event, userPrefs preference.UserPreferences) error {
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

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Notification\r\n\r\n%s", n.FromAddress, recipient, body.String()))

	cmd := exec.Command("sendmail", "-t")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		stdin.Write(msg)
	}()

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
