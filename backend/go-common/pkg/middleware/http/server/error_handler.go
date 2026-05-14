package middleware

import (
	"errors"
	"strings"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/commonerror"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		localizer := ctxmetadata.GetLocalizerFromContext(c.UserContext())
		lang := extractLang(c.Get("Accept-Language"))

		var commonErr *commonerror.CommonError
		if errors.As(err, &commonErr) {
			code := commonErrToHTTPStatus(commonErr)
			msg := localizeMsg(localizer, commonErr.Error())
			return c.Status(code).JSON(ErrorResponse{Error: msg})
		}

		if st, ok := status.FromError(err); ok && st.Code() != codes.OK {
			code := grpcCodeToHTTPStatus(st.Code())
			msg, fields := extractGRPCDetails(st, lang)
			return c.Status(code).JSON(ErrorResponse{Error: msg, Fields: fields})
		}

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			return c.Status(fiberErr.Code).JSON(ErrorResponse{Error: fiberErr.Message})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: localizeMsg(localizer, commonerror.ErrInternal.Error()),
		})
	}
}

func commonErrToHTTPStatus(err *commonerror.CommonError) int {
	switch err {
	case commonerror.ErrUnauthorized:
		return fiber.StatusUnauthorized
	case commonerror.ErrPermissionDenied:
		return fiber.StatusForbidden
	case commonerror.ErrNotFound:
		return fiber.StatusNotFound
	case commonerror.ErrAlreadyExists:
		return fiber.StatusConflict
	case commonerror.ErrInvalidArgument, commonerror.ErrFailedPrecondition, commonerror.ErrOutOfRange:
		return fiber.StatusBadRequest
	case commonerror.ErrResourceExhausted:
		return fiber.StatusTooManyRequests
	case commonerror.ErrDeadlineExceeded:
		return fiber.StatusGatewayTimeout
	case commonerror.ErrUnimplemented:
		return fiber.StatusNotImplemented
	case commonerror.ErrUnavailable:
		return fiber.StatusServiceUnavailable
	default:
		return fiber.StatusInternalServerError
	}
}

func grpcCodeToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.Canceled:
		return 499
	case codes.InvalidArgument, codes.FailedPrecondition, codes.OutOfRange:
		return fiber.StatusBadRequest
	case codes.DeadlineExceeded:
		return fiber.StatusGatewayTimeout
	case codes.NotFound:
		return fiber.StatusNotFound
	case codes.AlreadyExists, codes.Aborted:
		return fiber.StatusConflict
	case codes.PermissionDenied:
		return fiber.StatusForbidden
	case codes.ResourceExhausted:
		return fiber.StatusTooManyRequests
	case codes.Unimplemented:
		return fiber.StatusNotImplemented
	case codes.Unavailable:
		return fiber.StatusServiceUnavailable
	case codes.Unauthenticated:
		return fiber.StatusUnauthorized
	default:
		return fiber.StatusInternalServerError
	}
}

func extractGRPCDetails(st *status.Status, lang string) (string, map[string]string) {
	var mainMsg string
	var fields map[string]string

	prefix := lang + "#"
	for _, detail := range st.Details() {
		lm, ok := detail.(*errdetails.LocalizedMessage)
		if !ok {
			continue
		}
		if lm.Locale == lang {
			mainMsg = lm.Message
		} else if strings.HasPrefix(lm.Locale, prefix) {
			fieldName := lm.Locale[len(prefix):]
			if fields == nil {
				fields = make(map[string]string)
			}
			fields[fieldName] = lm.Message
		}
	}

	if mainMsg == "" {
		mainMsg = st.Message()
	}
	return mainMsg, fields
}

func localizeMsg(localizer *i18n.Localizer, id string) string {
	if localizer == nil {
		return id
	}
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:      id,
		DefaultMessage: &i18n.Message{ID: id, Other: id},
	})
	if err != nil || msg == "" {
		return id
	}
	return msg
}
