package senderservice

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/email-service/internal/config/settings"
	"github.com/ZaiiiRan/messenger/backend/email-service/internal/consumers/models"
	codemessage "github.com/ZaiiiRan/messenger/backend/email-service/internal/domain/code_message"
	"github.com/ZaiiiRan/messenger/backend/email-service/internal/transport/i18n"
	"github.com/ZaiiiRan/messenger/backend/email-service/internal/transport/smtp"
	"go.uber.org/zap"
)

type SenderService interface {
	SendCodeMessage(ctx context.Context, c models.CodeMessage, workerID string) error
}

type senderService struct {
	htmlGeneratorCfg *settings.HTMLGeneratorSettings
	smtpClient       *smtp.SMTPClient
	log              *zap.SugaredLogger
}

func NewSenderService(htmlGeneratorCfg settings.HTMLGeneratorSettings, smtpClient *smtp.SMTPClient, log *zap.SugaredLogger) SenderService {
	return &senderService{
		htmlGeneratorCfg: &htmlGeneratorCfg,
		smtpClient:       smtpClient,
		log:              log,
	}
}

func (s *senderService) SendCodeMessage(ctx context.Context, c models.CodeMessage, workerID string) error {
	l := s.log.With("op", "send_code_message", "worker_id", workerID)

	if c.ExpiresAt.Before(time.Now()) {
		l.Warnw("sender.send_code_message_failed.expired", "expires_at", c.ExpiresAt)
		return nil
	}

	codeMessage, err := codemessage.New(c.Id, c.Email, c.Code, c.LinkToken, c.CodeType, c.Language)
	if err != nil {
		l.Warnw("sender.send_code_message_failed.code_message_create_failed", "err", err)
		return err
	}

	loc := i18n.NewLocalizer(codeMessage.GetLanguage())
	err = codeMessage.GenerateHTML(s.htmlGeneratorCfg, loc)
	if err != nil {
		l.Warnw("sender.send_code_message_failed.generate_html_failed", "err", err)
		return err
	}

	err = s.smtpClient.SendHTMLMail(codeMessage.GetEmail(), codeMessage.GetSubject(loc), codeMessage.GetHTML())
	if err != nil {
		l.Errorw("sender.send_code_message_failed", "err", err)
		return err
	}

	l.Infow("sender.send_code_message.success")
	return nil
}
