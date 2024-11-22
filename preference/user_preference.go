package preference

type UserPreferences struct {
	UserID            string
	PreferredLanguage string
	EnabledChannels   map[string]bool
	EnabledEventTypes map[string]bool
}
