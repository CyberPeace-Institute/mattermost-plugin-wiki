package bot

import (
	pluginapi "github.com/mattermost/mattermost-plugin-api"
)

// Bot stores the information for the plugin configuration, and implements the Poster and Logger
// interfaces.
type Bot struct {
	pluginAPI  *pluginapi.Client
	botUserID  string
	logContext LogContext
}

// Logger interface - a logging system that will tee logs to a DM channel.
type Logger interface {
	With(LogContext) Logger
	Timed() Logger
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}

// New creates a new bot poster/logger.
func New(api *pluginapi.Client, botUserID string) *Bot {
	return &Bot{
		pluginAPI: api,
		botUserID: botUserID,
	}
}

// Clone shallow copies
func (b *Bot) clone() *Bot {
	return &Bot{
		pluginAPI:  b.pluginAPI,
		botUserID:  b.botUserID,
		logContext: b.logContext.copyShallow(),
	}
}
