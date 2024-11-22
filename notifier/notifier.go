package notifier

import (
	"github.com/pyxsoft/notificator/event"
	"github.com/pyxsoft/notificator/preference"
	"github.com/pyxsoft/notificator/template"
)

type Notifier interface {
	Send(event *event.Event, userPrefs preference.UserPreferences) error
	GetChannelName() string
	SetTemplateManager(tmplManager *template.TemplateManager)
}
