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

func New(username, channel string) *SlackHook {
	hook := &SlackHook{
		Channel: channel,
		Opts: []slack.MsgOption{
			slack.MsgOptionTS(time.Now().String()),
			slack.MsgOptionUsername(username),
		}}
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

func (hook *SlackHook) PostEphemeralAttachments(ctx context.Context, cli *slack.Client, userId string, attachments ...slack.Attachment) error {
	opts := append(hook.Opts, slack.MsgOptionAttachments(attachments...))
	 _, err := cli.PostEphemeralContext(ctx, hook.Channel, userId, opts...)
	if err != nil {
		return err
	}
	return nil
}


// GetChannelByName returns a `Channel` with the given name.
func (c *SlackHook) GetChannelByName(cli *slack.Client, name string) (*slack.Channel, error) {
	channels, err := cli.GetChannels(false)
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		if channel.Name == name {
			return &channel, nil
		}
	}
	return nil, nil
}

// GetChannelByName returns a `Channel` with the given name.
func (c *SlackHook) GetUserByEmail(ctx context.Context, cli *slack.Client, email string) (*slack.User, error) {
	return cli.GetUserByEmailContext(ctx, email)
}

type UserFunc func(u *slack.User) error


func (h *SlackHook) UserFunc(ctx context.Context, cli *slack.Client, email string, funcs ...UserFunc) error {
	usr, err := h.GetUserByEmail(ctx,cli,  email)
	if err != nil {
		return  err
	}
	for _, f := range funcs {
		if err := f(usr); err != nil {
			return err
		}
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

type EventHandler func(m slack.RTMEvent)

const (
	SLACK_COLOR_DANGER  string = "danger"
	SLACK_COLOR_WARNING string = "warning"
	SLACK_COLOR_GOOD    string = "good"
)


func EventLoop(cli *slack.Client, opts ...EventHandler) {
	rtm := cli.NewRTM()
	go rtm.ManageConnection()
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			for _, o := range opts {
				o(msg)
			}
		}
	}
}