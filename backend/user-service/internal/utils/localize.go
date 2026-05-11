package utils

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	appi18n "github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/i18n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func GetLocalizerFromContext(ctx context.Context) *i18n.Localizer {
	if loc, ok := ctx.Value(ctxmetadata.CtxKeyLocalizer{}).(*i18n.Localizer); ok && loc != nil {
		return loc
	}
	return CreateLocalizer(ctx)
}

func CreateLocalizer(ctx context.Context) *i18n.Localizer {
	return appi18n.NewLocalizer(ctxmetadata.GetLangFromIncomingContext(ctx))
}

func Localize(ctx context.Context, id string, fallback ...string) string {
	localizer := GetLocalizerFromContext(ctx)
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:      id,
		DefaultMessage: &i18n.Message{ID: id, Other: fallbackOrID(id, fallback)},
	})
	if err != nil {
		return fallbackOrID(id, fallback)
	}
	return msg
}

func fallbackOrID(id string, fallback []string) string {
	if len(fallback) > 0 && fallback[0] != "" {
		return fallback[0]
	}
	return id
}
