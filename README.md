# Notificator

**Notificator** is a flexible and extensible Go module for managing and sending notifications across multiple channels. It allows developers to easily configure, customize, and deliver notifications to users via email, Slack, Telegram, Mailchimp Transactional API, custom HTTP POST endpoints, and more.

## Features

- **Multi-Channel Support**: Send notifications through various channels, including:
    - Postfix (local email server)
    - Sendmail (local email server)
    - Mailchimp Transactional API
    - Telegram
    - Slack
    - Custom HTTP POST endpoints
    - Traditional SMTP email
- **Template-Based Notifications**: Create dynamic, localized notifications using text or HTML templates.
- **Event-Driven Architecture**: Register and handle different event types dynamically.
- **User Preferences**: Manage user-specific preferences for channels, languages, and enabled event types.
- **Localization**: Send notifications in the user's preferred language with fallback to default.
- **Customizability**: Easily add new notification channels and event types.

## Getting Started

### Installation

Install the module via `go get`:

```bash
go get github.com/pyxsoft/notificator
```

### Basic Usage

Here’s an example of how to set up and use **Notificator**:

#### 1. Import the module:

```go
import "github.com/pyxsoft/notificator"
```

#### 2. Define Templates

Create templates for each notification channel and event type in a structured directory. For example:

```
templates/
├── en/
│   ├── email/
│   │   └── user_registered.tmpl
│   ├── slack/
│   │   └── user_registered.tmpl
│   └── telegram/
│       └── user_registered.tmpl
```

#### 3. Initialize the Notificator System

Set up the system with your preferred channels and configurations:

```go
package main

import (
    "embed"
    "github.com/pyxsoft/notificator"
    "github.com/pyxsoft/notificator/notifier"
    "github.com/pyxsoft/notificator/preference"
    "github.com/pyxsoft/notificator/event"
)

//go:embed templates/*
var templatesFS embed.FS

func main() {
    // Configure your preferred notifiers
    notifiers := []notifier.Notifier{
        notifier.NewPostfixNotifier("no-reply@example.com"),
        notifier.NewSlackNotifier("https://hooks.slack.com/services/..."),
        notifier.NewTelegramNotifier("your-bot-token", "your-chat-id"),
    }

    // Initialize the notificator system
    system := notificator.NewNotificatorSystem(notificator.Config{
        TemplatesFS:       templatesFS,
        TemplatesBasePath: "templates",
        Notifiers:         notifiers,
    })

    // Register events
    system.EventRegistry.AddEvent("user_registered", map[string]interface{}{
        "Username": "",
        "Email":    "",
    })

    // Define user preferences
    userPrefs := preference.UserPreferences{
        UserID:            "user123",
        PreferredLanguage: "en",
        EnabledChannels: map[string]bool{
            "postfix":  true,
            "slack":    true,
            "telegram": true,
        },
        EnabledEventTypes: map[string]bool{
            "user_registered": true,
        },
    }

    // Create event data
    eventData := map[string]interface{}{
        "Username": "John Doe",
        "Email":    "john.doe@example.com",
    }

    // Trigger a notification
    userEvent := event.NewEvent("user_registered", eventData)
    system.NotificatorManager.SendToAll(userEvent, userPrefs)
}
```

### Supported Channels

- **Postfix**: Send emails via the local Postfix server in Linux.
- **Sendmail**: Send emails via the local Sendmail command in Linux.
- **Mailchimp Transactional API**: Deliver transactional emails through Mailchimp.
- **Telegram**: Notify users or groups using Telegram bots.
- **Slack**: Post messages to Slack channels via webhook URLs.
- **Custom HTTP POST**: Send JSON-encoded event data to a custom URL.
- **SMTP Email**: Use any SMTP server for email delivery.

### Key Benefits

- **Scalable and Modular**: Add new notification channels without altering core functionality.
- **Easy to Integrate**: Minimal setup to start delivering notifications.
- **Localized Templates**: Support multiple languages with flexible template management.
- **Event Customization**: Define and handle custom event types dynamically.

## Documentation

For detailed documentation, including setup, configuration, and advanced features, visit the [Notificator Documentation](https://github.com/pyxsoft/notificator/wiki).

## Contributing

Contributions are welcome! If you'd like to contribute, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Submit a pull request with a clear description of the changes.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
