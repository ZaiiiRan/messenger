package interfaces

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"
)

type EmailCodeTasksProducer interface {
	ProduceEmailCodeTask(ctx context.Context, email string, code *code.Code, language string) error
}