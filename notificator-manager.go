package notificator

import (
	"github.com/pyxsoft/notificator/event"
	"github.com/pyxsoft/notificator/notifier"
	"github.com/pyxsoft/notificator/preference"
)

type NotificatorManager struct {
	Notifiers     []notifier.Notifier
	EventRegistry *event.EventRegistry
}

func NewNotificatorManager(notifiers []notifier.Notifier, eventRegistry *event.EventRegistry) *NotificatorManager {
	return &NotificatorManager{
		Notifiers:     notifiers,
		EventRegistry: eventRegistry,
	}
}

func (m *NotificatorManager) SendToAll(event *event.Event, userPrefs preference.UserPreferences) {
	for _, n := range m.Notifiers {
		if !userPrefs.EnabledChannels[n.GetChannelName()] {
			continue
		}
		if !userPrefs.EnabledEventTypes[event.GetType()] {
			continue
		}
		if err := n.Send(event, userPrefs); err != nil {
			// Handle the error (e.g., log it)
		}
	}
}
