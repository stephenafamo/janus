package authboss

import (
	"context"

	"github.com/stephenafamo/mailer"
	"github.com/volatiletech/authboss/v3"
)

type Mailer struct {
	Mailer mailer.Mailer
}

func (a Mailer) Send(ctx context.Context, abEmail authboss.Email) error {
	email := mailer.Email{
		Subject: abEmail.Subject,

		To:      abEmail.To,
		ToNames: abEmail.ToNames,

		Cc:      abEmail.Cc,
		CcNames: abEmail.CcNames,

		Bcc:      abEmail.Bcc,
		BccNames: abEmail.BccNames,

		From:     abEmail.From,
		FromName: abEmail.FromName,

		ReplyTo:     abEmail.ReplyTo,
		ReplyToName: abEmail.ReplyToName,

		TextBody: abEmail.TextBody,
		HTMLBody: abEmail.HTMLBody,
	}

	_, _, err := a.Mailer.Send(ctx, email)
	return err
}
