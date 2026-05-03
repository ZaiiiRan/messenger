package ctxmetadata

import (
	"context"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"google.golang.org/grpc/metadata"
)

type CtxKeyLocalizer struct{}

const AcceptLanguageKey = "x-accept-language"

func WithLocalizer(ctx context.Context, localizer *i18n.Localizer) context.Context {
	return context.WithValue(ctx, CtxKeyLocalizer{}, localizer)
}

func GetLangFromIncomingContext(ctx context.Context) string {
	lang := "en"
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(AcceptLanguageKey); len(values) > 0 && values[0] != "" {
			lang = values[0]
		}
	}
	return lang
}

func ForwardLangToOutgoingContext(ctx context.Context) context.Context {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if langs := md.Get(AcceptLanguageKey); len(langs) > 0 && langs[0] != "" {
			return metadata.AppendToOutgoingContext(ctx, AcceptLanguageKey, langs[0])
		}
	}
	return ctx
}
