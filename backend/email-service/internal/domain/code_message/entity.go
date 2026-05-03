package codemessage

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ZaiiiRan/messenger/backend/email-service/internal/config/settings"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type CodeType string

const (
	CodeTypeActivation    CodeType = "activation"
	CodeTypePasswordReset CodeType = "password_reset"
)

type CodeMessage struct {
	id        int64
	email     string
	code      string
	linkToken string
	language  string
	codeType  CodeType
	html      string
}

func New(
	id int64,
	email string,
	code string,
	linkToken string,
	codeType string,
	language string,
) (*CodeMessage, error) {
	ct, err := toCodeType(codeType)
	if err != nil {
		return nil, err
	}

	return &CodeMessage{
		id:        id,
		email:     email,
		code:      code,
		linkToken: linkToken,
		language:  language,
		codeType:  ct,
	}, nil
}

func (m *CodeMessage) GetID() int64          { return m.id }
func (m *CodeMessage) GetEmail() string      { return m.email }
func (m *CodeMessage) GetHTML() string       { return m.html }
func (m *CodeMessage) GetCodeType() CodeType { return m.codeType }
func (m *CodeMessage) GetLanguage() string   { return m.language }

func (m *CodeMessage) GenerateHTML(cfg *settings.HTMLGeneratorSettings, localizer *i18n.Localizer) error {
	switch m.codeType {
	case CodeTypeActivation:
		return m.generateActivationHTML(cfg.BaseUrlForActivation, localizer)
	case CodeTypePasswordReset:
		return m.generatePasswordResetHTML(cfg.BaseUrlForPasswordReset, localizer)
	default:
		return fmt.Errorf("unsupported code type: %s", m.codeType)
	}
}

func (m *CodeMessage) GetSubject(localizer *i18n.Localizer) string {
	switch m.codeType {
	case CodeTypeActivation:
		return loc(localizer, "email.activation.subject")
	case CodeTypePasswordReset:
		return loc(localizer, "email.password_reset.subject")
	default:
		return ""
	}
}

func (m *CodeMessage) generateActivationHTML(baseURL string, localizer *i18n.Localizer) error {
	r := strings.NewReplacer(
		"{SUBTITLE}", loc(localizer, "email.activation.subtitle"),
		"{BODY}", loc(localizer, "email.activation.body"),
		"{CODE_LABEL}", loc(localizer, "email.code_label"),
		"{CODE}", m.code,
		"{DIVIDER}", loc(localizer, "email.divider"),
		"{TOKEN_URL}", baseURL+"?token="+url.QueryEscape(m.linkToken),
		"{BUTTON}", loc(localizer, "email.activation.button"),
		"{FOOTER}", loc(localizer, "email.activation.footer"),
	)
	m.html = r.Replace(activationHTMLTpl)
	return nil
}

func (m *CodeMessage) generatePasswordResetHTML(baseURL string, localizer *i18n.Localizer) error {
	r := strings.NewReplacer(
		"{SUBTITLE}", loc(localizer, "email.password_reset.subtitle"),
		"{WARNING}", loc(localizer, "email.password_reset.warning"),
		"{BODY}", loc(localizer, "email.password_reset.body"),
		"{CODE_LABEL}", loc(localizer, "email.code_label"),
		"{CODE}", m.code,
		"{DIVIDER}", loc(localizer, "email.divider"),
		"{TOKEN_URL}", baseURL+"?token="+url.QueryEscape(m.linkToken),
		"{BUTTON}", loc(localizer, "email.password_reset.button"),
		"{FOOTER}", loc(localizer, "email.password_reset.footer"),
	)
	m.html = r.Replace(passwordResetHTMLTpl)
	return nil
}

func loc(localizer *i18n.Localizer, id string) string {
	s, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: id})
	if err != nil {
		return id
	}
	return s
}

func toCodeType(codeType string) (CodeType, error) {
	switch codeType {
	case "activation":
		return CodeTypeActivation, nil
	case "password_reset":
		return CodeTypePasswordReset, nil
	default:
		return "", fmt.Errorf("invalid code type: %s", codeType)
	}
}
