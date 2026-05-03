package smtp

import (
	"fmt"
	"sync"
	"time"

	"github.com/ZaiiiRan/messenger/backend/email-service/internal/config/settings"
	"gopkg.in/gomail.v2"
)

type SMTPClient struct {
	dialer *gomail.Dialer
	sender gomail.SendCloser

	cfg *settings.SMTPClientSettings

	mu sync.Mutex
}

func New(cfg settings.SMTPClientSettings) *SMTPClient {
	d := gomail.NewDialer(
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
	)

	return &SMTPClient{
		dialer: d,
		cfg:    &cfg,
	}
}

func (c *SMTPClient) SendHTMLMail(to, subject, body string) error {
	m := gomail.NewMessage()
	from := c.cfg.From
	if from == "" {
		from = c.cfg.Username
	}
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return c.sendWithRetry(m)
}

func (c *SMTPClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closeConn()
}

func (c *SMTPClient) sendWithRetry(m *gomail.Message) error {
	delay := time.Millisecond * time.Duration(c.cfg.MaxRetries)

	for attempt := 1; attempt <= int(c.cfg.MaxRetries); attempt++ {
		err := c.trySend(m)
		if err == nil {
			return nil
		}

		if attempt == int(c.cfg.MaxRetries) {
			return fmt.Errorf("failed to send email after %d attempts: %w", attempt, err)
		}

		time.Sleep(delay)

		d := c.cfg.RetryDelayMS * 2
		if d < c.cfg.RetryMaxDelayMS {
			delay = time.Millisecond * time.Duration(d)
		} else {
			delay = time.Millisecond * time.Duration(c.cfg.RetryMaxDelayMS)
		}
	}

	return nil
}

func (c *SMTPClient) trySend(m *gomail.Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.ensureConnected(); err != nil {
		return err
	}

	if err := gomail.Send(c.sender, m); err != nil {
		_ = c.closeConn()
		return err
	}

	return nil
}

func (c *SMTPClient) ensureConnected() error {
	if c.sender != nil {
		return nil
	}

	s, err := c.dialer.Dial()
	if err != nil {
		return fmt.Errorf("failed to dial SMTP server: %w", err)
	}

	c.sender = s
	return nil
}

func (c *SMTPClient) closeConn() error {
	if c.sender == nil {
		return nil
	}
	err := c.sender.Close()
	c.sender = nil
	return err
}
