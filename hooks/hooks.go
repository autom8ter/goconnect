package hooks

import (
	"context"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	"time"
)

type SlackHook struct {
	Channel string `validate:"required"`
	Opts    []slack.MsgOption
}

func New(username, channel string, broadcast bool) *SlackHook {
	hook := &SlackHook{
		Channel: channel,
		Opts: []slack.MsgOption{
			slack.MsgOptionTS(time.Now().String()),
			slack.MsgOptionUsername(username),
		}}

	if broadcast {
		hook.Opts = append(hook.Opts, slack.MsgOptionBroadcast())
	}
	return hook
}

func (hook *SlackHook) PostAttachments(ctx context.Context, cli *slack.Client, attachments ...slack.Attachment) error {
	opts := append(hook.Opts, slack.MsgOptionAttachments(attachments...))
	_, _, err := cli.PostMessageContext(ctx, hook.Channel, opts...)
	if err != nil {
		return err
	}
	return nil
}

func (hook *SlackHook) PostLogEntry(ctx context.Context, cli *slack.Client, author, icon, title string, sourceEntry *logrus.Entry) error {
	var messageFields []slack.AttachmentField
	for key, value := range sourceEntry.Data {
		message := slack.AttachmentField{
			Title: key,
			Value: value.(string),
			Short: true,
		}

		messageFields = append(messageFields, message)
	}
	attachment := slack.Attachment{
		Color:      getColor(sourceEntry.Level),
		AuthorName: author,
		AuthorIcon: icon,
		Title:      title,
		Text:       sourceEntry.Message,
		Fields:     messageFields,
	}
	return hook.PostAttachments(ctx, cli, attachment)
}

func getColor(level logrus.Level) string {
	switch level {
	case logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel:
		return SLACK_COLOR_DANGER
	case logrus.WarnLevel:
		return SLACK_COLOR_WARNING
	case logrus.InfoLevel:
		return SLACK_COLOR_GOOD
	default:
		return ""
	}
}

const (
	SLACK_COLOR_DANGER  string = "danger"
	SLACK_COLOR_WARNING string = "warning"
	SLACK_COLOR_GOOD    string = "good"
)
