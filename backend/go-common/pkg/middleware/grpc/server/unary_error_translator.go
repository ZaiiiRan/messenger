package middleware

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func ErrorTranslatorMiddleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		localizer := ctxmetadata.GetLocalizerFromContext(ctx)
		if localizer == nil {
			return nil, err
		}

		locale := ctxmetadata.GetLangFromIncomingContext(ctx)
		return nil, translateError(localizer, locale, err)
	}
}

func translateError(localizer *i18n.Localizer, locale string, err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	originalID := st.Message()
	newSt := status.New(st.Code(), originalID)

	localizedMsg := &errdetails.LocalizedMessage{
		Locale:  locale,
		Message: localizeID(localizer, originalID),
	}

	badReq := extractTranslatedBadRequest(st, localizer, locale)
	if badReq != nil {
		newSt, _ = newSt.WithDetails(localizedMsg, badReq)
	} else {
		newSt, _ = newSt.WithDetails(localizedMsg)
	}

	return newSt.Err()
}

func extractTranslatedBadRequest(st *status.Status, localizer *i18n.Localizer, locale string) *errdetails.BadRequest {
	for _, detail := range st.Details() {
		d, ok := detail.(*errdetails.BadRequest)
		if !ok {
			continue
		}
		violations := make([]*errdetails.BadRequest_FieldViolation, len(d.FieldViolations))
		for i, fv := range d.FieldViolations {
			violations[i] = &errdetails.BadRequest_FieldViolation{
				Field:       fv.Field,
				Description: fv.Description,
				LocalizedMessage: &errdetails.LocalizedMessage{
					Locale:  locale,
					Message: localizeID(localizer, fv.Description),
				},
			}
		}
		return &errdetails.BadRequest{FieldViolations: violations}
	}
	return nil
}

func localizeID(localizer *i18n.Localizer, id string) string {
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:      id,
		DefaultMessage: &i18n.Message{ID: id, Other: id},
	})
	if err != nil || msg == "" {
		return id
	}
	return msg
}
