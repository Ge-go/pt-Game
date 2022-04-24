package email

import "context"

var DefaultSender *EmailSender

func NewEmailSender(domain, uri, gameId, sigKey string, channelId, source int) (*EmailSender, error) {
	sender := newEmailSender(domain, uri, gameId, sigKey, channelId, source)
	DefaultSender = sender

	return sender, nil
}

func SendEmail(ctx context.Context, email, subject, template string) error {
	return DefaultSender.SendEmail(ctx, email, subject, template)
}
